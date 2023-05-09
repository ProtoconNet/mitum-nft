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

var nftSignItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(NFTSignItemProcessor)
	},
}

var nftSignProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(NFTSignProcessor)
	},
}

func (NFTSign) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type NFTSignItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   NFTSignItem
}

func (ipp *NFTSignItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	nid := ipp.item.NFT()

	st, err := existsState(NFTStateKey(ipp.item.contract, ipp.item.symbol, CollectionKey), "key of design", getStateFunc)
	if err != nil {
		return errors.Errorf("collection design not found, %q: %w", nid.Collection(), err)
	}

	design, err := StateCollectionValue(st)
	if err != nil {
		return errors.Errorf("collection design value not found, %q: %w", ipp.item.symbol, err)
	}

	if !design.Active() {
		return errors.Errorf("deactivated collection, %q", nid.Collection())
	}
	st, err = existsState(extensioncurrency.StateKeyContractAccount(ipp.item.contract), "contract account", getStateFunc)
	if err != nil {
		return errors.Errorf("parent not found, %q: %w", design.Parent(), err)
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return errors.Errorf("contract account value not found, %q: %w", ipp.item.contract, err)
	}

	if !ca.IsActive() {
		return errors.Errorf("deactivated contract account, %q", ipp.item.contract)
	}

	st, err = existsState(StateKeyNFT(ipp.item.contract, ipp.item.symbol, nid), "key of nft", getStateFunc)
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

	switch ipp.item.Qualification() {
	case CreatorQualification:
		if nv.Creators().IsSignedByAddress(ipp.sender) {
			return errors.Errorf("already signed nft, %q-%q", ipp.sender, nv.ID())
		}
	case CopyrighterQualification:
		if nv.Copyrighters().IsSignedByAddress(ipp.sender) {
			return errors.Errorf("already signed nft, %q-%q", ipp.sender, nv.ID())
		}
	default:
		return errors.Errorf("wrong qualification, %q", ipp.item.Qualification())
	}

	return nil
}

func (ipp *NFTSignItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := existsState(StateKeyNFT(ipp.item.contract, ipp.item.symbol, nid), "key of nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %q: %w", nid, err)
	}

	nv, err := StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %q: %w", nid, err)
	}

	var signers nft.Signers

	switch ipp.item.Qualification() {
	case CreatorQualification:
		signers = nv.Creators()
	case CopyrighterQualification:
		signers = nv.Copyrighters()
	default:
		return nil, errors.Errorf("wrong qualification, %q", ipp.item.Qualification())
	}

	idx := signers.IndexByAddress(ipp.sender)
	if idx < 0 {
		return nil, errors.Errorf("not signer of nft, %q-%q", ipp.sender, nv.ID())
	}

	signer := nft.NewSigner(signers.Signers()[idx].Account(), signers.Signers()[idx].Share(), true)
	if err := signer.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid signer, %q", signer.Account())
	}

	sns := &signers
	if err := sns.SetSigner(signer); err != nil {
		return nil, errors.Errorf("failed to set signer for signers, %q: %w", signer, err)
	}

	var n nft.NFT
	if ipp.item.Qualification() == CreatorQualification {
		n = nft.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), nv.Approved(), *sns, nv.Copyrighters())
	} else {
		n = nft.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), nv.Approved(), nv.Creators(), *sns)
	}

	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %q: %w", n.ID(), err)
	}

	sts := make([]base.StateMergeValue, 1)

	sts[0] = NewStateMergeValue(StateKeyNFT(ipp.item.contract, ipp.item.symbol, n.ID()), NewNFTStateValue(n))

	return sts, nil
}

func (ipp *NFTSignItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = NFTSignItem{}
	nftSignItemProcessorPool.Put(ipp)

	return nil
}

type NFTSignProcessor struct {
	*base.BaseOperationProcessor
}

func NewNFTSignProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new NFTSignProcessor")

		nopp := nftSignProcessorPool.Get()
		opp, ok := nopp.(*NFTSignProcessor)
		if !ok {
			return nil, e(nil, "expected NFTSignProcessor, not %T", nopp)
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

func (opp *NFTSignProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess NFTSign")

	fact, ok := op.Fact().(NFTSignFact)
	if !ok {
		return ctx, nil, e(nil, "expected NFTSignFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot sign nfts, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := nftSignItemProcessorPool.Get()
		ipc, ok := ip.(*NFTSignItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected NFTSignItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess NFTSignItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *NFTSignProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process NFTSign")

	fact, ok := op.Fact().(NFTSignFact)
	if !ok {
		return nil, nil, e(nil, "expected NFTSignFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue

	for _, item := range fact.Items() {
		ip := nftSignItemProcessorPool.Get()
		ipc, ok := ip.(*NFTSignItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected NFTSignItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process MintItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	fitems := fact.Items()
	items := make([]CollectionItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := CalculateCollectionItemsFee(getStateFunc, items)
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

func (opp *NFTSignProcessor) Close() error {
	nftSignProcessorPool.Put(opp)

	return nil
}
