package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum2/base"
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

func (AddSignature) Process(
	ctx context.Context, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type SignItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	item   AddSignatureItem
}

func (ipp *SignItemProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	e := util.StringError("preprocess SignItemProcessor")

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

	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.NFTStateKey(ipp.item.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection state for contract account %v: %v", it.Contract(), err))
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Errorf("nft collection state value for contract account %v: %v", it.Contract(), err))

	}

	if !design.Active() {
		return e.Wrap(
			errors.Errorf("nft collection in contract account %v has already been deactivated", ipp.item.Contract()))
	}

	st, err = state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Errorf("nft idx %v in contract account %v", nid, ipp.item.Contract()))
	}

	nv, err := statenft.StateNFTValue(st)
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

func (ipp *SignItemProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %v: %w", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
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

	sts := make([]mitumbase.StateMergeValue, 1)

	sts[0] = state.NewStateMergeValue(statenft.StateKeyNFT(ipp.item.Contract(), n.ID()), statenft.NewNFTStateValue(n))

	return sts, nil
}

func (ipp *SignItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AddSignatureItem{}
	signItemProcessorPool.Put(ipp)

	return
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

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender %v is contract account", fact.Sender(), cErr)), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := signItemProcessorPool.Get()
		ipc, ok := ip.(*SignItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected SignItemProcessor, not %T", ip)), nil
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

func (opp *SignProcessor) Process( // nolint:dupl
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Sign")

	fact, ok := op.Fact().(AddSignatureFact)
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

func (opp *SignProcessor) Close() error {
	signProcessorPool.Put(opp)

	return nil
}
