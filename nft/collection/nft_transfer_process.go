package collection

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var nftTransferItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(NFTTransferItemProcessor)
	},
}

var nftTransferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(NFTTransferProcessor)
	},
}

func (NFTTransfer) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type NFTTransferItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   NFTTransferItem
}

func (ipp *NFTTransferItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	receiver := ipp.item.Receiver()

	if err := checkExistsState(currency.StateKeyAccount(receiver), getStateFunc); err != nil {
		return errors.Errorf("receiver not found, %q: %w", receiver, err)
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(receiver), getStateFunc); err != nil {
		return errors.Errorf("contract account cannot receive nfts, %q: %w", receiver, err)
	}

	nid := ipp.item.NFT()

	st, err := existsState(NFTStateKey(ipp.item.contract, ipp.item.collection, CollectionKey), "design", getStateFunc)
	if err != nil {
		return errors.Errorf("collection design not found, %q: %w", nid.Collection(), err)
	}

	design, err := StateCollectionValue(st)
	if err != nil {
		return errors.Errorf("collection design not found, %q: %w", nid.Collection(), err)
	}
	if !design.Active() {
		return errors.Errorf("deactivated collection, %q", design.Collection())
	}

	st, err = existsState(extensioncurrency.StateKeyContractAccount(design.Parent()), "key of contract account", getStateFunc)
	if err != nil {
		return errors.Errorf("parent not found, %q: %w", design.Parent(), err)
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return errors.Errorf("parent account value not found, %q: %w", design.Parent(), err)
	}

	if !ca.IsActive() {
		return errors.Errorf("deactivated contract account, %q", design.Parent())
	}

	st, err = existsState(StateKeyNFT(ipp.item.contract, ipp.item.collection, nid), "key of nft", getStateFunc)
	if err != nil {
		return errors.Errorf("nft not found, %q: %w", nid, err)
	}

	nv, err := StateNFTValue(st)
	if err != nil {
		return errors.Errorf("nft value not found, %q: %w", nid, err)
	}

	if !nv.Active() {
		return errors.Errorf("burned nft, %q", nid)
	}

	if !(nv.Owner().Equal(ipp.sender) || nv.Approved().Equal(ipp.sender)) {
		if st, err := existsState(StateKeyOperators(ipp.item.contract, ipp.item.collection, nv.Owner()), "operators", getStateFunc); err != nil {
			return errors.Errorf("unauthorized sender, %q: %w", ipp.sender, err)
		} else if box, err := StateOperatorsBookValue(st); err != nil {
			return errors.Errorf("operator book value not found, %q: %w", ipp.sender, err)
		} else if !box.Exists(ipp.sender) {
			return errors.Errorf("unauthorized sender, %q", ipp.sender)
		}
	}

	return nil
}

func (ipp *NFTTransferItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	receiver := ipp.item.Receiver()
	nid := ipp.item.NFT()

	st, err := existsState(StateKeyNFT(ipp.item.contract, ipp.item.collection, nid), "key of nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %q: %w", nid, err)
	}

	nv, err := StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %q: %w", nid, err)
	}

	n := nft.NewNFT(nid, nv.Active(), receiver, nv.NFTHash(), nv.URI(), receiver, nv.Creators(), nv.Copyrighters())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %q: %w", nid, err)
	}

	sts := make([]base.StateMergeValue, 1)

	sts[0] = NewStateMergeValue(StateKeyNFT(ipp.item.contract, ipp.item.collection, ipp.item.NFT()), NewNFTStateValue(n))

	return sts, nil
}

func (ipp *NFTTransferItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = NFTTransferItem{}

	nftTransferItemProcessorPool.Put(ipp)

	return nil
}

type NFTTransferProcessor struct {
	*base.BaseOperationProcessor
}

func NewNFTTransferProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new NFTTransferProcessor")

		nopp := nftTransferProcessorPool.Get()
		opp, ok := nopp.(*NFTTransferProcessor)
		if !ok {
			return nil, e(nil, "expected NFTTransferProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *NFTTransferProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess NFTTransfer")

	fact, ok := op.Fact().(NFTTransferFact)
	if !ok {
		return ctx, nil, e(nil, "expected NFTTransferFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot transfer nfts, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := nftTransferItemProcessorPool.Get()
		ipc, ok := ip.(*NFTTransferItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected NFTTransferItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess NFTTransferItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *NFTTransferProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process NFTTransfer")

	fact, ok := op.Fact().(NFTTransferFact)
	if !ok {
		return nil, nil, e(nil, "expected NFTTransferFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := nftTransferItemProcessorPool.Get()
		ipc, ok := ip.(*NFTTransferItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected NFTTransferItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process NFTTransferItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	required, err := opp.calculateItemsFee(op, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *NFTTransferProcessor) Close() error {
	nftTransferProcessorPool.Put(opp)

	return nil
}

func (opp *NFTTransferProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := op.Fact().(NFTTransferFact)
	if !ok {
		return nil, errors.Errorf("expected NFTTransferFact, not %T", op.Fact())
	}

	items := make([]CollectionItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateCollectionItemsFee(getStateFunc, items)
}
