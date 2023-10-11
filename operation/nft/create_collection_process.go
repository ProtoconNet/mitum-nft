package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var createCollectionProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateCollectionProcessor)
	},
}

func (CreateCollection) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateCollectionProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewCreateCollectionProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateCollectionProcessor")

		nopp := createCollectionProcessorPool.Get()
		opp, ok := nopp.(*CreateCollectionProcessor)
		if !ok {
			return nil, errors.Errorf("expected CreateCollectionProcessor, not %T", nopp)
		}

		b, err := mitumbase.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *CreateCollectionProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CreateCollection")

	fact, ok := op.Fact().(CreateCollectionFact)
	if !ok {
		return ctx, nil, e.Errorf("expected CreateCollectionFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	_, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender is contract account, %q", fact.Sender()), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}

	if !(ca.Owner().Equal(fact.sender) || ca.IsOperator(fact.Sender())) {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender is neither the owner nor the operator of the target contract account, %q", fact.sender), nil
	}

	if ca.IsActive() {
		return nil, mitumbase.NewBaseOperationProcessReasonError("a contract account is already used, %q", fact.Contract().String()), nil
	}

	if err := currencystate.CheckNotExistsState(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("collection design already exists, %q; %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckNotExistsState(statenft.NFTStateKey(fact.contract, statenft.LastIDXKey), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("last index of collection design already exists, %q; %w", fact.Contract(), err), nil
	}

	whitelist := fact.WhiteList()
	for _, white := range whitelist {
		if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(white), getStateFunc); err != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError("whitelist account not found, %q; %w", white, err), nil
		} else if err = currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(white), getStateFunc); err != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError("whitelist account is contract account, %q; %w", white, err), nil
		}
	}

	return ctx, nil, nil
}

func (opp *CreateCollectionProcessor) Process(
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CreateCollection")

	fact, ok := op.Fact().(CreateCollectionFact)
	if !ok {
		return nil, nil, e.Errorf("expected CreateCollectionFact, not %T", op.Fact())
	}

	sts := make([]mitumbase.StateMergeValue, 4)

	policy := types.NewCollectionPolicy(fact.Name(), fact.Royalty(), fact.URI(), fact.WhiteList())
	design := types.NewDesign(fact.Contract(), fact.Sender(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("invalid collection design, %q; %w", fact.Contract(), err), nil
	}

	sts[0] = currencystate.NewStateMergeValue(
		statenft.NFTStateKey(design.Parent(), statenft.CollectionKey),
		statenft.NewCollectionStateValue(design),
	)
	sts[1] = currencystate.NewStateMergeValue(
		statenft.NFTStateKey(design.Parent(), statenft.LastIDXKey),
		statenft.NewLastNFTIndexStateValue(0),
	)

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}
	nca := ca.SetIsActive(true)

	sts[2] = currencystate.NewStateMergeValue(
		stateextension.StateKeyContractAccount(fact.Contract()),
		stateextension.NewContractAccountStateValue(nca),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to check fee of currency, %q; %w", fact.Currency(), err), nil
	}

	st, err = currencystate.ExistsState(statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender balance not found, %q; %w", fact.Sender(), err), nil
	}
	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := statecurrency.StateBalanceValue(st); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get balance value, %q; %w", statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[3] = currencystate.NewStateMergeValue(
		sb.Key(),
		statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CreateCollectionProcessor) Close() error {
	createCollectionProcessorPool.Put(opp)

	return nil
}
