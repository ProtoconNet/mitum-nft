package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
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
	ctx context.Context, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	item   TransferItem
}

func (ipp *TransferItemProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	err := state.CheckExistsState(stateextension.StateKeyContractAccount(ipp.item.Contract()), getStateFunc)
	if err != nil {
		return errors.Errorf("contract account not found, %q; %w", ipp.item.Contract(), err)
	}

	receiver := ipp.item.Receiver()

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(receiver), getStateFunc); err != nil {
		return errors.Errorf("receiver not found, %q; %w", receiver, err)
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(receiver), getStateFunc); err != nil {
		return errors.Errorf("contract account cannot receive nfts, %q; %w", receiver, err)
	}

	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.NFTStateKey(ipp.item.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return errors.Errorf("collection design not found, %q; %w", ipp.item.Contract(), err)
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return errors.Errorf("collection design not found, %q; %w", ipp.item.Contract(), err)
	}
	if !design.Active() {
		return errors.Errorf("deactivated collection, %q", ipp.item.Contract())
	}

	st, err = state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "key of nft", getStateFunc)
	if err != nil {
		return errors.Errorf("nft not found, %q; %w", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return errors.Errorf("nft value not found, %q; %w", nid, err)
	}

	if !nv.Active() {
		return errors.Errorf("burned nft, %q", nid)
	}

	if !(nv.Owner().Equal(ipp.sender) || nv.Approved().Equal(ipp.sender)) {
		if st, err := state.ExistsState(statenft.StateKeyOperators(ipp.item.Contract(), nv.Owner()), "operators", getStateFunc); err != nil {
			return errors.Errorf("unauthorized sender, %q; %w", ipp.sender, err)
		} else if box, err := statenft.StateOperatorsBookValue(st); err != nil {
			return errors.Errorf("operator book value not found, %q; %w", ipp.sender, err)
		} else if !box.Exists(ipp.sender) {
			return errors.Errorf("unauthorized sender, %q", ipp.sender)
		}
	}

	return nil
}

func (ipp *TransferItemProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	receiver := ipp.item.Receiver()
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "key of nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %q; %w", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %q; %w", nid, err)
	}

	n := types.NewNFT(nid, nv.Active(), receiver, nv.NFTHash(), nv.URI(), receiver, nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %q; %w", nid, err)
	}

	sts := make([]mitumbase.StateMergeValue, 1)

	sts[0] = state.NewStateMergeValue(statenft.StateKeyNFT(ipp.item.Contract(), ipp.item.NFT()), statenft.NewNFTStateValue(n))

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
	*mitumbase.BaseOperationProcessor
}

func NewTransferProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new TransferProcessor")

		nopp := transferProcessorPool.Get()
		opp, ok := nopp.(*TransferProcessor)
		if !ok {
			return nil, e.Errorf("expected TransferProcessor, not %T", nopp)
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

func (opp *TransferProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Transfer")

	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, nil, e.Errorf("expected TransferFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("contract account cannot transfer nfts, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := transferItemProcessorPool.Get()
		ipc, ok := ip.(*TransferItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected TransferItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("fail to preprocess TransferItem; %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *TransferProcessor) Process( // nolint:dupl
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Transfer")

	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return nil, nil, e.Errorf("expected TransferFact, not %T", op.Fact())
	}

	var sts []mitumbase.StateMergeValue // nolint:prealloc
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
			return nil, mitumbase.NewBaseOperationProcessReasonError("failed to process TransferItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	items := make([]CollectionItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiverBalSts, required, err := CalculateCollectionItemsFee(getStateFunc, items)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.Sender(), required, getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for cid := range sb {
		v, ok := sb[cid].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		_, feeReceiverFound := feeReceiverBalSts[cid]

		if feeReceiverFound && (sb[cid].Key() != feeReceiverBalSts[cid].Key()) {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiverBalSts[cid].Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, mitumbase.NewBaseOperationProcessReasonError("expected %T, not %T", statecurrency.BalanceStateValue{}, feeReceiverBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiverBalSts[cid].Key(),
					statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
						return statecurrency.NewBalanceStateValueMerger(height, feeReceiverBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *TransferProcessor) Close() error {
	transferProcessorPool.Put(opp)

	return nil
}
