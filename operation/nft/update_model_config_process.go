package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var updateModelConfigProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateModelConfigProcessor)
	},
}

func (UpdateModelConfig) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateModelConfigProcessor struct {
	*base.BaseOperationProcessor
}

func NewUpdateModelConfigProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new UpdateModelConfigProcessor")

		nopp := updateModelConfigProcessorPool.Get()
		opp, ok := nopp.(*UpdateModelConfigProcessor)
		if !ok {
			return nil, errors.Errorf("expected UpdateModelConfigProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *UpdateModelConfigProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(UpdateModelConfigFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", UpdateModelConfigFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	whitelist := fact.Whitelist()
	for _, white := range whitelist {
		if _, _, _, cErr := cstate.ExistsCAccount(white, "whitelist", true, false, getStateFunc); cErr != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
					Errorf("%v: whitelist %v is contract account", cErr, white)), nil
		}
	}

	st, err := cstate.ExistsState(state.NFTStateKey(fact.Contract(), state.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft service state for contract account %v", fact.Contract())), nil

	}

	design, err := state.StateCollectionValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft service state value for contract account %v", fact.Contract())), nil
	}

	if !design.Active() {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("nft service in contract account %v has already been deactivated ", fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateModelConfigProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(UpdateModelConfigFact)

	st, err := cstate.ExistsState(state.NFTStateKey(fact.contract, state.CollectionKey), "design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design not found, %v: %w", fact.Contract(), err), nil
	}

	design, err := state.StateCollectionValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("collection design value not found, %v: %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue
	whitelist := fact.Whitelist()
	for _, white := range whitelist {
		smv, err := cstate.CreateNotExistAccount(white, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
		} else if smv != nil {
			sts = append(sts, smv)
		}
	}

	de := types.NewDesign(
		design.Contract(), design.Creator(), design.Active(), design.Count(),
		types.NewCollectionPolicy(fact.name, fact.royalty, fact.uri, fact.whitelist),
	)
	sts = append(sts, cstate.NewStateMergeValue(state.NFTStateKey(fact.contract, state.CollectionKey), state.NewCollectionStateValue(de)))

	return sts, nil, nil
}

func (opp *UpdateModelConfigProcessor) Close() error {
	updateModelConfigProcessorPool.Put(opp)

	return nil
}
