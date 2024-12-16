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
	"github.com/ProtoconNet/mitum2/base"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var updateModelConfigProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateModelConfigProcessor)
	},
}

func (UpdateModelConfig) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateModelConfigProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewUpdateModelConfigProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new UpdateModelConfigProcessor")

		nopp := updateModelConfigProcessorPool.Get()
		opp, ok := nopp.(*UpdateModelConfigProcessor)
		if !ok {
			return nil, errors.Errorf("expected UpdateModelConfigProcessor, not %T", nopp)
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

func (opp *UpdateModelConfigProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(UpdateModelConfigFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", UpdateModelConfigFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := state.CheckExistsState(currency.DesignStateKey(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id, %v", fact.Currency())), nil
	}

	whitelist := fact.Whitelist()
	for _, white := range whitelist {
		if _, _, _, cErr := state.ExistsCAccount(white, "whitelist", true, false, getStateFunc); cErr != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
					Errorf("%v: whitelist %v is contract account", cErr, white)), nil
		}
	}

	st, err := state.ExistsState(statenft.NFTStateKey(fact.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft collection state for contract account %v", fact.Contract())), nil

	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft collection state value for contract account %v", fact.Contract())), nil
	}

	if !design.Active() {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("collection in contract account %v has already been deactivated ", fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateModelConfigProcessor) Process(
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(UpdateModelConfigFact)

	st, err := state.ExistsState(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design not found, %v: %w", fact.Contract(), err), nil
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("collection design value not found, %v: %w", fact.Contract(), err), nil
	}

	var sts []mitumbase.StateMergeValue
	whitelist := fact.Whitelist()
	for _, white := range whitelist {
		smv, err := state.CreateNotExistAccount(white, getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("%w", err), nil
		} else if smv != nil {
			sts = append(sts, smv)
		}
	}

	de := types.NewDesign(
		design.Contract(), design.Creator(), design.Active(),
		types.NewCollectionPolicy(fact.name, fact.royalty, fact.uri, fact.whitelist),
	)
	sts = append(sts, state.NewStateMergeValue(statenft.NFTStateKey(fact.contract, statenft.CollectionKey), statenft.NewCollectionStateValue(de)))

	return sts, nil, nil
}

func (opp *UpdateModelConfigProcessor) Close() error {
	updateModelConfigProcessorPool.Put(opp)

	return nil
}
