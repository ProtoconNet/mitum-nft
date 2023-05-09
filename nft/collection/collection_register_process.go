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

var collectionRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionRegisterProcessor)
	},
}

func (CollectionRegister) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CollectionRegisterProcessor struct {
	*base.BaseOperationProcessor
}

func NewCollectionRegisterProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CollectionRegisterProcessor")

		nopp := collectionRegisterProcessorPool.Get()
		opp, ok := nopp.(*CollectionRegisterProcessor)
		if !ok {
			return nil, errors.Errorf("expected CollectionRegisterProcessor, not %T", nopp)
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

func (opp *CollectionRegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess CollectionRegister")

	fact, ok := op.Fact().(CollectionRegisterFact)
	if !ok {
		return ctx, nil, e(nil, "expected CollectionRegisterFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender is contract account, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	st, err := existsState(extensioncurrency.StateKeyContractAccount(fact.Form().Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("target contract account not found, %q: %w", fact.Form().Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q: %w", fact.Form().Contract(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return ctx, base.NewBaseOperationProcessReasonError("sender is not owner of contract account, %q, %q", fact.Sender(), ca.Owner()), nil
	}

	if !ca.IsActive() {
		return ctx, base.NewBaseOperationProcessReasonError("deactivated contract account, %q", fact.Form().Contract()), nil
	}

	if err := checkNotExistsState(NFTStateKey(fact.Form().contract, fact.Form().Symbol(), CollectionKey), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("collection design already exists, %q: %w", fact.Form().Symbol(), err), nil
	}

	if err := checkNotExistsState(NFTStateKey(fact.Form().contract, fact.Form().Symbol(), LastIDXKey), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("last index of collection design already exists, %q: %w", fact.Form().Symbol(), err), nil
	}

	whites := fact.Form().Whites()
	for _, white := range whites {
		if err := checkExistsState(currency.StateKeyAccount(white), getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError("whitelist account not found, %q: %w", white, err), nil
		} else if err = checkNotExistsState(extensioncurrency.StateKeyContractAccount(white), getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError("whitelist account is contract account, %q: %w", white, err), nil
		}
	}

	return ctx, nil, nil
}

func (opp *CollectionRegisterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process CollectionRegister")

	fact, ok := op.Fact().(CollectionRegisterFact)
	if !ok {
		return nil, nil, e(nil, "expected CollectionRegisterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 3)

	policy := NewCollectionPolicy(fact.Form().Name(), fact.Form().Royalty(), fact.Form().URI(), fact.Form().Whites())
	design := NewCollectionDesign(fact.Form().Contract(), fact.Sender(), fact.Form().Symbol(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid collection design, %q: %w", fact.Form().Symbol(), err), nil
	}

	sts[0] = NewStateMergeValue(
		NFTStateKey(design.Parent(), design.Symbol(), CollectionKey),
		NewCollectionStateValue(design),
	)
	sts[1] = NewStateMergeValue(
		NFTStateKey(design.Parent(), design.Symbol(), LastIDXKey),
		NewLastNFTIndexStateValue(0),
	)

	currencyPolicy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err := existsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
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
	sts[2] = currency.NewBalanceStateMergeValue(
		sb.Key(),
		currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CollectionRegisterProcessor) Close() error {
	collectionRegisterProcessorPool.Put(opp)

	return nil
}
