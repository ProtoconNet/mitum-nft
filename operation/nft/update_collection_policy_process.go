package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var updateCollectionPolicyProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateCollectionPolicyProcessor)
	},
}

func (UpdateCollectionPolicy) Process(
	ctx context.Context, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateCollectionPolicyProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewUpdateCollectionPolicyProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new UpdateCollectionPolicyProcessor")

		nopp := updateCollectionPolicyProcessorPool.Get()
		opp, ok := nopp.(*UpdateCollectionPolicyProcessor)
		if !ok {
			return nil, errors.Errorf("expected UpdateCollectionPolicyProcessor, not %T", nopp)
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

func (opp *UpdateCollectionPolicyProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess UpdateCollectionPolicy")
	fact, ok := op.Fact().(UpdateCollectionPolicyFact)
	if !ok {
		return ctx, nil, e.Errorf("not UpdateCollectionPolicyFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("contract account cannot update collection policy, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	st, err := state.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("contract account value not found, %q; %w", fact.Contract(), err), nil
	}

	if !(ca.Owner().Equal(fact.Sender()) || ca.IsOperator(fact.Sender())) {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender is neither the owner nor the operator of the target contract account, %q", fact.Sender()), nil
	}

	st, err = state.ExistsState(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), "key of collection design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design not found, %q; %w", fact.Contract(), err), nil
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design value not found, %q; %w", fact.Contract(), err), nil
	}

	if !design.Active() {
		return nil, mitumbase.NewBaseOperationProcessReasonError("deactivated collection, %q", fact.Contract()), nil
	}

	if !design.Parent().Equal(fact.Contract()) {
		return nil, mitumbase.NewBaseOperationProcessReasonError("parent not found, %q; %w", design.Parent(), err), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateCollectionPolicyProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process UpdateCollectionPolicy")
	fact, ok := op.Fact().(UpdateCollectionPolicyFact)
	if !ok {
		return nil, nil, e.Errorf("expected UpdateCollectionPolicyFact, not %T", op.Fact())
	}

	st, err := state.ExistsState(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), "key of design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design not found, %q; %w", fact.Contract(), err), nil
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design value not found, %q; %w", fact.Contract(), err), nil
	}

	var sts []mitumbase.StateMergeValue

	de := types.NewDesign(
		design.Parent(),
		design.Creator(),
		design.Active(),
		types.NewCollectionPolicy(fact.name, fact.royalty, fact.uri, fact.whitelist),
	)
	sts = append(sts, state.NewStateMergeValue(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), statenft.NewCollectionStateValue(de)))

	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	if currencyPolicy.Feeer().Receiver() == nil {
		return sts, nil, nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	senderBalSt, err := state.ExistsState(
		statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
		"key of sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"sender balance not found, %q; %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := statecurrency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
		return nil, nil, err
	} else if feeRcvrSt, found, err := getStateFunc(statecurrency.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
		return nil, nil, err
	} else if !found {
		return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
	} else if feeRcvrSt.Key() != senderBalSt.Key() {
		r, ok := feeRcvrSt.Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, errors.Errorf("expected %T, not %T", statecurrency.BalanceStateValue{}, feeRcvrSt.Value())
		}
		sts = append(sts, common.NewBaseStateMergeValue(
			feeRcvrSt.Key(),
			statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
			func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
				return statecurrency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
			},
		))

		sts = append(sts, common.NewBaseStateMergeValue(
			senderBalSt.Key(),
			statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
			func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
				return statecurrency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
			},
		))
	}

	return sts, nil, nil
}

func (opp *UpdateCollectionPolicyProcessor) Close() error {
	updateCollectionPolicyProcessorPool.Put(opp)

	return nil
}
