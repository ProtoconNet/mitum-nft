package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
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
	e := util.StringError("preprocess TransferItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id, %v", it.Currency())))
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(it.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(cErr)
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(it.Receiver(), "receiver", true, false, getStateFunc); aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(cErr))
	}

	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.NFTStateKey(ipp.item.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection, %s: %v", it.Contract(), err))
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection, %s: %v", it.Contract(), err))
	}
	if !design.Active() {
		return e.Wrap(errors.Errorf("deactivated collection, %v", ipp.item.Contract()))
	}

	st, err = state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Errorf("nft, %v: %v", nid, err))
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return e.Wrap(common.ErrStateValInvalid.Errorf("nft, %v: %v", nid, err))
	}

	if !nv.Active() {
		return e.Wrap(errors.Errorf("burned nft, %v", nid))
	}

	if !(nv.Owner().Equal(ipp.sender) || nv.Approved().Equal(ipp.sender)) {
		if st, err := state.ExistsState(statenft.StateKeyOperators(ipp.item.Contract(), nv.Owner()), "operators", getStateFunc); err != nil {
			return e.Wrap(common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v: %v", ipp.sender, err)))
		} else if box, err := statenft.StateOperatorsBookValue(st); err != nil {
			return e.Wrap(common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v: %v", ipp.sender, err)))
		} else if !box.Exists(ipp.sender) {
			return e.Wrap(common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v", ipp.sender)))
		}
	}

	if it.receiver.Equal(nv.Owner()) {
		return e.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver, %v is same with nft owner", it.receiver)))
	}

	if nv.Owner().Equal(ipp.sender) && ipp.sender.Equal(it.receiver) {
		return e.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver, %v is same with sender", ipp.sender)))
	}

	return nil
}

func (ipp *TransferItemProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	receiver := ipp.item.Receiver()
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %v: %w", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %v: %w", nid, err)
	}

	n := types.NewNFT(nid, nv.Active(), receiver, nv.NFTHash(), nv.URI(), receiver, nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %v: %w", nid, err)
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
	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender account is contract account, %v", fact.Sender(), cErr)), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := transferItemProcessorPool.Get()
		ipc, ok := ip.(*TransferItemProcessor)
		if !ok {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected TransferItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
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
