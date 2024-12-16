package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ccstate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   ApproveItem
}

func (ipp *ApproveItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess ApproveItemProcessor")

	if err := cstate.CheckExistsState(
		ccstate.DesignStateKey(ipp.item.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id, %v", ipp.item.Currency())))
	}

	st, err := cstate.ExistsState(
		state.NFTStateKey(ipp.item.Contract(), state.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(
			common.ErrServiceNF.Wrap(
				errors.Errorf("nft collection state for contract account %v", ipp.item.Contract())))
	}

	design, err := state.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(
			common.ErrServiceNF.Wrap(
				errors.Errorf("nft collection state value for contract account %v", ipp.item.Contract())))
	}

	if !design.Active() {
		return e.Wrap(
			common.ErrServiceNF.Wrap(
				errors.Errorf("collection in the contract account %v has been deactived", ipp.item.Contract())))
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		ipp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: approved %v is contract account", cErr, ipp.item.Approved())))
	}

	st, err = cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), ipp.item.nftIdx), "nft", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Wrap(
			errors.Errorf("nft idx %v in contract account %v", ipp.item.nftIdx, ipp.item.Contract())))
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return e.Wrap(common.ErrStateValInvalid.Wrap(
			errors.Errorf("nft idx %v in contract account %v", ipp.item.nftIdx, ipp.item.Contract())))
	}

	if !nv.Active() {
		return e.Wrap(common.ErrStateValInvalid.Wrap(
			errors.Errorf("burned nft idx %v in contract account %v", ipp.item.nftIdx, ipp.item.Contract())))
	}

	if ipp.item.Approved().Equal(nv.Approved()) {
		return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("already approved %v", ipp.item.Approved())))
	}

	if !nv.Owner().Equal(ipp.sender) {
		if err := cstate.CheckExistsState(ccstate.AccountStateKey(nv.Owner()), getStateFunc); err != nil {
			return e.Wrap(
				common.ErrAccountNF.Wrap(errors.Errorf("nft owner %v for nft idx %v", nv.Owner(), ipp.item.nftIdx)))
		}

		st, err = cstate.ExistsState(
			state.StateKeyOperators(ipp.item.Contract(), ipp.sender), "operators", getStateFunc)
		if err != nil {
			return e.Wrap(
				common.ErrStateNF.Wrap(
					common.ErrAccountNAth.Wrap(
						errors.Errorf(
							"sender %v neither nft owner nor operator for nft idx %v",
							ipp.sender, ipp.item.nftIdx))))
		}

		operators, err := state.StateOperatorsBookValue(st)
		if err != nil {
			return e.Wrap(
				common.ErrStateValInvalid.Wrap(
					common.ErrAccountNAth.Wrap(
						errors.Errorf(
							"sender %v neither nft owner nor operator for nft idx %v",
							ipp.sender, ipp.item.nftIdx))))
		}

		if !operators.Exists(ipp.sender) {
			return e.Wrap(
				common.ErrAccountNAth.Wrap(
					errors.Errorf(
						"sender %v neither nft owner nor operator for nft idx %v: sender is not in operators",
						ipp.sender, ipp.item.nftIdx)))
		}
	}

	return nil
}

func (ipp *ApproveItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(ipp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	nid := ipp.item.NFTIdx()

	st, err := cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, util.ErrNotFound.Errorf("nft, %v: %v", nid, err)
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return nil, util.ErrNotFound.Errorf("nft value, %v: %v", nid, err)
	}

	n := types.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), ipp.item.Approved(), nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, err
	}

	sts = append(sts, cstate.NewStateMergeValue(st.Key(), state.NewNFTStateValue(n)))

	return sts, nil
}

func (ipp *ApproveItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveItem{}

	approveItemProcessorPool.Put(ipp)

	return
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new ApproveProcessor")

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
		if !ok {
			return nil, e.Errorf("expected ApproveProcessor, not %T", nopp)
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ApproveFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected ApproveItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Approve")

	fact, _ := op.Fact().(ApproveFact)

	var sts []base.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected ApproveItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("process ApproveItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)

	return nil
}
