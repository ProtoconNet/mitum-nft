package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	statee "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var registerModelProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterModelProcessor)
	},
}

func (RegisterModel) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterModelProcessor struct {
	*base.BaseOperationProcessor
}

func NewRegisterModelProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RegisterModelProcessor")

		nopp := registerModelProcessorPool.Get()
		opp, ok := nopp.(*RegisterModelProcessor)
		if !ok {
			return nil, errors.Errorf("expected RegisterModelProcessor, not %T", nopp)
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

func (opp *RegisterModelProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", RegisterModelFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if found, _ := cstate.CheckNotExistsState(state.NFTStateKey(fact.contract, state.CollectionKey), getStateFunc); found {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceE).Errorf("nft collection for contract account %v", fact.Contract())), nil
	}

	if found, _ := cstate.CheckNotExistsState(state.NFTStateKey(fact.contract, state.LastIDXKey), getStateFunc); found {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceE).Errorf("nft collection for contract account %v: last index already exists", fact.Contract())), nil
	}

	whitelist := fact.WhiteList()
	for _, white := range whitelist {
		if _, _, _, cErr := cstate.ExistsCAccount(white, "whitelist", true, false, getStateFunc); cErr != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
					Errorf("%v: whitelist %v is contract account", cErr, white)), nil
		}
	}

	return ctx, nil, nil
}

func (opp *RegisterModelProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(RegisterModelFact)
	var sts []base.StateMergeValue
	whitelist := fact.WhiteList()
	for _, white := range whitelist {
		smv, err := cstate.CreateNotExistAccount(white, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
		} else if smv != nil {
			sts = append(sts, smv)
		}
	}

	policy := types.NewCollectionPolicy(fact.Name(), fact.Royalty(), fact.URI(), fact.WhiteList())
	design := types.NewDesign(fact.Contract(), fact.Sender(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid collection design, %v: %w", fact.Contract(), err), nil
	}

	sts = append(sts, cstate.NewStateMergeValue(
		state.NFTStateKey(design.Contract(), state.CollectionKey),
		state.NewCollectionStateValue(design),
	))
	sts = append(sts, cstate.NewStateMergeValue(
		state.NFTStateKey(design.Contract(), state.LastIDXKey),
		state.NewLastNFTIndexStateValue(0),
	))

	st, _ := cstate.ExistsState(statee.StateKeyContractAccount(fact.Contract()), "contract account", getStateFunc)
	ca, _ := statee.StateContractAccountValue(st)
	nca := ca.SetIsActive(true)

	sts = append(sts, cstate.NewStateMergeValue(
		statee.StateKeyContractAccount(fact.Contract()),
		statee.NewContractAccountStateValue(nca),
	))

	return sts, nil, nil
}

func (opp *RegisterModelProcessor) Close() error {
	registerModelProcessorPool.Put(opp)

	return nil
}
