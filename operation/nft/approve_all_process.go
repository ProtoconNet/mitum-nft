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

var delegateItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateItemProcessor)
	},
}

var delegateProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateProcessor)
	},
}

func (ApproveAll) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DelegateItemProcessor struct {
	h      util.Hash
	sender base.Address
	box    *types.AllApprovedBook
	item   ApproveAllItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess DelegateItemProcessor")

	if _, _, _, cErr := cstate.ExistsCAccount(
		ipp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: approved %v is contract account", cErr, ipp.item.Approved())))
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(ipp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	if ipp.box == nil {
		return nil, errors.Errorf(
			"nft box not found, %v", state.StateKeyOperators(ipp.item.Contract(), ipp.sender))
	}

	switch ipp.item.Mode() {
	case ApproveAllAllow:
		if err := ipp.box.Append(ipp.item.Approved()); err != nil {
			return nil, err
		}
	case ApproveAllCancel:
		if err := ipp.box.Remove(ipp.item.Approved()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("wrong mode for delegate item, %v: \"allow\" | \"cancel\"", ipp.item.Mode())
	}

	ipp.box.Sort(true)

	return sts, nil
}

func (ipp *DelegateItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveAllItem{}
	ipp.box = nil

	delegateItemProcessorPool.Put(ipp)

	return
}

type DelegateProcessor struct {
	*base.BaseOperationProcessor
}

func NewDelegateProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new DelegateProcessor")

		nopp := delegateProcessorPool.Get()
		opp, ok := nopp.(*DelegateProcessor)
		if !ok {
			return nil, e.Errorf("expected DelegateProcessor, not %T", nopp)
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

func (opp *DelegateProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApproveAllFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ApproveAllFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		st, err := cstate.ExistsState(
			state.NFTStateKey(item.Contract(), state.CollectionKey), "design", getStateFunc)
		if err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("nft collection, %v: %v", item.Contract(), err)), nil
		}

		design, err := state.StateCollectionValue(st)
		if err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("nft collection, %v: %v", item.Contract(), err)), nil
		}

		if !design.Active() {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("collection in contract account %v has already been deactivated",
						item.Contract())), nil
		}
	}

	for _, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected DelegateItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *DelegateProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Delegate")

	fact, _ := op.Fact().(ApproveAllFact)
	boxes := map[string]*types.AllApprovedBook{}
	for _, item := range fact.Items() {
		ak := state.StateKeyOperators(item.contract, fact.Sender())

		var operators types.AllApprovedBook
		switch st, found, err := getStateFunc(ak); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to get state of operators book, %v: %w", ak, err), nil
		case !found:
			operators = types.NewAllApprovedBook(nil)
		default:
			o, err := state.StateOperatorsBookValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					"operators book value not found, %v: %w", ak, err), nil
			} else {
				operators = *o
			}
		}
		boxes[ak] = &operators
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*DelegateItemProcessor, len(fact.items))
	for i, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected %T, not %T", &DelegateItemProcessor{}, ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = boxes[state.StateKeyOperators(item.contract, fact.Sender())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process DelegateItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for ak, box := range boxes {
		bv := cstate.NewStateMergeValue(ak, state.NewOperatorsBookStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	return sts, nil, nil
}

func (opp *DelegateProcessor) Close() error {
	delegateProcessorPool.Put(opp)

	return nil
}
