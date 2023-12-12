package nft

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var signItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignItemProcessor)
	},
}

var signProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignProcessor)
	},
}

func (Sign) Process(
	ctx context.Context, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type SignItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	item   SignItem
}

func (ipp *SignItemProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.NFTStateKey(ipp.item.contract, statenft.CollectionKey), "key of design", getStateFunc)
	if err != nil {
		return errors.Errorf("collection design not found, %q; %w", ipp.item.contract, err)
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return errors.Errorf("collection design value not found, %q; %w", ipp.item.contract, err)
	}

	if !design.Active() {
		return errors.Errorf("deactivated collection, %q", ipp.item.contract)
	}
	st, err = state.ExistsState(stateextension.StateKeyContractAccount(ipp.item.contract), "contract account", getStateFunc)
	if err != nil {
		return errors.Errorf("parent not found, %q; %w", design.Parent(), err)
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return errors.Errorf("contract account value not found, %q; %w", ipp.item.contract, err)
	}

	if !ca.IsActive() {
		return errors.Errorf("deactivated contract account, %q", ipp.item.contract)
	}

	st, err = state.ExistsState(statenft.StateKeyNFT(ipp.item.contract, nid), "key of nft", getStateFunc)
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

	if nv.Creators().IsSignedByAddress(ipp.sender) {
		return errors.Errorf("already signed nft, %q-%q", ipp.sender, nv.ID())
	}

	return nil
}

func (ipp *SignItemProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.StateKeyNFT(ipp.item.contract, nid), "key of nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %q; %w", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %q; %w", nid, err)
	}

	signers := nv.Creators()

	idx := signers.IndexByAddress(ipp.sender)
	if idx < 0 {
		return nil, errors.Errorf("not signer of nft, %q-%q", ipp.sender, nv.ID())
	}

	signer := types.NewSigner(signers.Signers()[idx].Account(), signers.Signers()[idx].Share(), true)
	if err := signer.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid signer, %q", signer.Account())
	}

	sns := &signers
	if err := sns.SetSigner(signer); err != nil {
		return nil, errors.Errorf("failed to set signer for signers, %q; %w", signer, err)
	}

	n := types.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), nv.Approved(), *sns)

	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %q; %w", n.ID(), err)
	}

	sts := make([]mitumbase.StateMergeValue, 1)

	sts[0] = state.NewStateMergeValue(statenft.StateKeyNFT(ipp.item.contract, n.ID()), statenft.NewNFTStateValue(n))

	return sts, nil
}

func (ipp *SignItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = SignItem{}
	signItemProcessorPool.Put(ipp)

	return nil
}

type SignProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewSignProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new SignProcessor")

		nopp := signProcessorPool.Get()
		opp, ok := nopp.(*SignProcessor)
		if !ok {
			return nil, e.Errorf("expected SignProcessor, not %T", nopp)
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

func (opp *SignProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Sign")

	fact, ok := op.Fact().(SignFact)
	if !ok {
		return ctx, nil, e.Errorf("expected SignFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("contract account cannot sign nfts, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := signItemProcessorPool.Get()
		ipc, ok := ip.(*SignItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected SignItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("fail to preprocess SignItem; %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *SignProcessor) Process( // nolint:dupl
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Sign")

	fact, ok := op.Fact().(SignFact)
	if !ok {
		return nil, nil, e.Errorf("expected SignFact, not %T", op.Fact())
	}

	var sts []mitumbase.StateMergeValue

	for _, item := range fact.Items() {
		ip := signItemProcessorPool.Get()
		ipc, ok := ip.(*SignItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected SignItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("failed to process MintItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	items := make([]CollectionItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiveBalSts, required, err := CalculateCollectionItemsFee(getStateFunc, items)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for cid := range sb {
		v, ok := sb[cid].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		if sb[cid].Key() != feeReceiveBalSts[cid].Key() {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiveBalSts[cid].Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, mitumbase.NewBaseOperationProcessReasonError("expected %T, not %T", statecurrency.BalanceStateValue{}, feeReceiveBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiveBalSts[cid].Key(),
					statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
						return statecurrency.NewBalanceStateValueMerger(height, feeReceiveBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *SignProcessor) Close() error {
	signProcessorPool.Put(opp)

	return nil
}
