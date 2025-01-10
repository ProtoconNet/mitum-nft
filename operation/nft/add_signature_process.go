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

var AddSignatureItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddSignatureItemProcessor)
	},
}

var AddSignatureProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddSignatureProcessor)
	},
}

func (AddSignature) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AddSignatureItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   AddSignatureItem
}

func (ipp *AddSignatureItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess AddSignatureItemProcessor")

	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	nid := ipp.item.NFT()

	st, err := cstate.ExistsState(state.NFTStateKey(ipp.item.Contract(), state.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection state for contract account %v: %v", it.Contract(), err))
	}

	design, err := state.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection state value for contract account %v: %v", it.Contract(), err))

	}

	if !design.Active() {
		return e.Wrap(
			errors.Errorf("nft collection in contract account %v has already been deactivated", ipp.item.Contract()))
	}

	st, err = cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Errorf("nft idx %v in contract account %v", nid, ipp.item.Contract()))
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return e.Wrap(
			common.ErrStateValInvalid.Errorf("nft idx %v in contract account %v", nid, ipp.item.Contract()))
	}

	if !nv.Active() {
		return e.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("burned nft idx %v in contract account %v", nid, ipp.item.Contract())))
	}

	if nv.Creators().IsSignedByAddress(ipp.sender) {
		return e.Wrap(errors.Errorf("already signed nft idx %v by creator %v", nv.ID(), ipp.sender))
	}

	return nil
}

func (ipp *AddSignatureItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := cstate.ExistsState(state.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %v: %w", nid, err)
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %v: %w", nid, err)
	}

	signers := nv.Creators()

	idx := signers.IndexByAddress(ipp.sender)
	if idx < 0 {
		return nil, errors.Errorf("not signer of nft, %v-%v", ipp.sender, nv.ID())
	}

	signer := types.NewSigner(signers.Signers()[idx].Address(), signers.Signers()[idx].Share(), true)
	if err := signer.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid signer, %v", signer.Address())
	}

	sns := &signers
	if err := sns.SetSigner(signer); err != nil {
		return nil, errors.Errorf("failed to set signer for signers, %v: %w", signer, err)
	}

	n := types.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), nv.Approved(), *sns)

	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %v: %w", n.ID(), err)
	}

	sts := make([]base.StateMergeValue, 1)

	sts[0] = cstate.NewStateMergeValue(state.StateKeyNFT(ipp.item.Contract(), n.ID()), state.NewNFTStateValue(n))

	return sts, nil
}

func (ipp *AddSignatureItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AddSignatureItem{}
	AddSignatureItemProcessorPool.Put(ipp)

	return
}

type AddSignatureProcessor struct {
	*base.BaseOperationProcessor
}

func NewSignProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AddSignatureProcessor")

		nopp := AddSignatureProcessorPool.Get()
		opp, ok := nopp.(*AddSignatureProcessor)
		if !ok {
			return nil, e.Errorf("expected SignProcessor, not %T", nopp)
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

func (opp *AddSignatureProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(AddSignatureFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", AddSignatureFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := AddSignatureItemProcessorPool.Get()
		ipc, ok := ip.(*AddSignatureItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected AddSignatureItemProcessor, not %T", ip)), nil
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

func (opp *AddSignatureProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process AddSignature")

	fact, _ := op.Fact().(AddSignatureFact)
	var sts []base.StateMergeValue

	for _, item := range fact.Items() {
		ip := AddSignatureItemProcessorPool.Get()
		ipc, ok := ip.(*AddSignatureItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected AddSignatureItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AddSignatureItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	return sts, nil, nil
}

func (opp *AddSignatureProcessor) Close() error {
	AddSignatureProcessorPool.Put(opp)

	return nil
}
