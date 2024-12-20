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

var transferItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferItemProcessor)
	},
}

var transferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferProcessor)
	},
}

func (Transfer) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   TransferItem
}

func (ipp *TransferItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess TransferItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := cstate.CheckExistsState(ccstate.DesignStateKey(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", it.Currency())))
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		it.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: receiver %v is contract account", cErr, it.Receiver())))
	}

	nid := ipp.item.NFT()

	st, err := cstate.ExistsState(
		state.NFTStateKey(ipp.item.Contract(), state.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(
			common.ErrStateNF.Wrap(
				common.ErrServiceNF.Errorf("nft collection state for contract account %v", it.Contract())))
	}

	design, err := state.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(
			common.ErrStateValInvalid.Wrap(
				common.ErrServiceNF.Errorf("nft collection state value for contract account %v", it.Contract())))
	}
	if !design.Active() {
		return e.Wrap(
			errors.Errorf(
				"nft collection in contract account %v has already been deactivated ", ipp.item.Contract()))
	}

	st, err = cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Errorf("nft idx %v in contract account %v", nid, it.Contract()))
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return e.Wrap(common.ErrStateValInvalid.Errorf("nft idx %v in contract account %v", nid, it.Contract()))
	}

	if !nv.Active() {
		return e.Wrap(errors.Errorf("burned nft idx %v in contract account %v", nid, ipp.item.Contract()))
	}

	if !(nv.Owner().Equal(ipp.sender) || nv.Approved().Equal(ipp.sender)) {
		if st, err := cstate.ExistsState(
			state.StateKeyOperators(ipp.item.Contract(), nv.Owner()), "operators", getStateFunc); err != nil {
			return e.Wrap(
				common.ErrAccountNAth.Wrap(
					errors.Errorf(
						"sender %v neither nft owner nor operator for nft idx %v in contract account %v: operators state not found",
						ipp.sender, nid, ipp.item.Contract())))

		} else if box, err := state.StateOperatorsBookValue(st); err != nil {
			return e.Wrap(
				common.ErrAccountNAth.Wrap(
					errors.Errorf("sender %v neither nft owner nor operator for nft idx %v in contract account %v: operators state value not found",
						ipp.sender, nid, ipp.item.Contract())))
		} else if !box.Exists(ipp.sender) {
			return e.Wrap(common.ErrValueInvalid.Wrap(
				common.ErrAccountNAth.Wrap(
					errors.Errorf("sender %v neither nft owner nor operator for nft idx %v in contract account %v: sender is not in operators ",
						ipp.sender, nid, ipp.item.Contract()))))
		}
	}

	if it.receiver.Equal(nv.Owner()) {
		return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("receiver %v is same with nft owner", it.receiver)))
	}

	if nv.Owner().Equal(ipp.sender) && ipp.sender.Equal(it.receiver) {
		return e.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with sender", it.receiver)))
	}

	return nil
}

func (ipp *TransferItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	receiver := ipp.item.Receiver()
	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(receiver, getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	nid := ipp.item.NFT()

	st, err := cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %v: %v", nid, err)
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %v: %v", nid, err)
	}

	n := types.NewNFT(nid, nv.Active(), receiver, nv.NFTHash(), nv.URI(), receiver, nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %v: %v", nid, err)
	}

	sts = append(
		sts,
		cstate.NewStateMergeValue(
			state.StateKeyNFT(ipp.item.Contract(), ipp.item.NFT()), state.NewNFTStateValue(n)),
	)

	return sts, nil
}

func (ipp *TransferItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = TransferItem{}

	transferItemProcessorPool.Put(ipp)

	return
}

type TransferProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransferProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new TransferProcessor")

		nopp := transferProcessorPool.Get()
		opp, ok := nopp.(*TransferProcessor)
		if !ok {
			return nil, e.Errorf("expected TransferProcessor, not %T", nopp)
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

func (opp *TransferProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := transferItemProcessorPool.Get()
		ipc, ok := ip.(*TransferItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected TransferItemProcessor, not %T", ip)), nil
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

func (opp *TransferProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Transfer")

	fact, _ := op.Fact().(TransferFact)
	var sts []base.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := transferItemProcessorPool.Get()
		ipc, ok := ip.(*TransferItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected TransferItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process TransferItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	return sts, nil, nil
}

func (opp *TransferProcessor) Close() error {
	transferProcessorPool.Put(opp)

	return nil
}
