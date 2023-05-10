package collection

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var collectionPolicyUpdaterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionPolicyUpdaterProcessor)
	},
}

func (CollectionPolicyUpdater) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CollectionPolicyUpdaterProcessor struct {
	*base.BaseOperationProcessor
}

func NewCollectionPolicyUpdaterProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CollectionPolicyUpdaterProcessor")

		nopp := collectionPolicyUpdaterProcessorPool.Get()
		opp, ok := nopp.(*CollectionPolicyUpdaterProcessor)
		if !ok {
			return nil, errors.Errorf("expected CollectionPolicyUpdaterProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *CollectionPolicyUpdaterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess CollectionPolicyUpdater")

	fact, ok := op.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return ctx, nil, e(nil, "not CollectionPolicyUpdaterFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot update collection policy, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	st, err := existsState(NFTStateKey(fact.contract, fact.collection, CollectionKey), "key of collection design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design not found, %q: %w", fact.Collection(), err), nil
	}

	design, err := StateCollectionValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design value not found, %q: %w", fact.Collection(), err), nil
	}

	if !design.Active() {
		return nil, base.NewBaseOperationProcessReasonError("deactivated collection, %q", fact.Collection()), nil
	}

	if !design.Creator().Equal(fact.Sender()) {
		return nil, base.NewBaseOperationProcessReasonError("not creator of collection design, %q", fact.Collection()), nil
	}

	st, err = existsState(extensioncurrency.StateKeyContractAccount(design.Parent()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("parent not found, %q: %w", design.Parent(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %q: %w", design.Parent(), err), nil
	}

	if !ca.IsActive() {
		return nil, base.NewBaseOperationProcessReasonError("deactivated contract account, %q", design.Parent()), nil
	}

	return ctx, nil, nil
}

func (opp *CollectionPolicyUpdaterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process CollectionPolicyUpdater")

	fact, ok := op.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return nil, nil, e(nil, "expected CollectionPolicyUpdaterFact, not %T", op.Fact())
	}

	st, err := existsState(NFTStateKey(fact.contract, fact.collection, CollectionKey), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design not found, %q: %w", fact.Collection(), err), nil
	}

	design, err := StateCollectionValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design value not found, %q: %w", fact.Collection(), err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	de := NewCollectionDesign(design.Parent(), design.Creator(), design.Collection(), design.Active(), fact.Policy())
	sts[0] = NewStateMergeValue(NFTStateKey(fact.contract, design.Collection(), CollectionKey), NewCollectionStateValue(de))

	currencyPolicy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = existsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := currency.NewBalanceStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = currency.NewBalanceStateMergeValue(
		sb.Key(),
		currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CollectionPolicyUpdaterProcessor) Close() error {
	collectionPolicyUpdaterProcessorPool.Put(opp)

	return nil
}
