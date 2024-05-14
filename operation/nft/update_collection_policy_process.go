package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"

	"github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/base"
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
	_ context.Context, _ mitumbase.GetStateFunc,
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
	fact, ok := op.Fact().(UpdateCollectionPolicyFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", UpdateCollectionPolicyFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := state.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id, %v", fact.Currency())), nil
	}

	if _, _, aErr, cErr := state.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	whitelist := fact.Whitelist()
	for _, white := range whitelist {
		if _, _, aErr, cErr := state.ExistsCAccount(white, "whitelist", true, false, getStateFunc); aErr != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", aErr)), nil
		} else if cErr != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
					Errorf("%v: whitelist account is contract account, %v", white, cErr)), nil
		}
	}

	_, cSt, aErr, cErr := state.ExistsCAccount(fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	_, err := stateextension.CheckCAAuthFromState(cSt, fact.Sender())
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	st, err := state.ExistsState(statenft.NFTStateKey(fact.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft collection, %s: %v", fact.Contract(), err)), nil

	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft collection, %s: %v", fact.Contract(), err)), nil
	}

	if !design.Active() {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("deactivated collection, %v", fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateCollectionPolicyProcessor) Process(
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process UpdateCollectionPolicy")
	fact, ok := op.Fact().(UpdateCollectionPolicyFact)
	if !ok {
		return nil, nil, e.Errorf("expected UpdateCollectionPolicyFact, not %T", op.Fact())
	}

	st, err := state.ExistsState(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design not found, %v: %w", fact.Contract(), err), nil
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design value not found, %v: %w", fact.Contract(), err), nil
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
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %v: %w", fact.Currency(), err), nil
	}

	if currencyPolicy.Feeer().Receiver() == nil {
		return sts, nil, nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %v: %w",
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
			"sender balance not found, %v: %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := statecurrency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get balance value, %v: %w",
			statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %v",
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
