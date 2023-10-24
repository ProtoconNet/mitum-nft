package nft

import (
	"context"
	"sync"

	statenft "github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var mintItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintItemProcessor)
	},
}

var mintProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintProcessor)
	},
}

func (Mint) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type MintItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   MintItem
	idx    uint64
	box    *types.NFTBox
}

func (ipp *MintItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	if err := state.CheckExistsState(stateextension.StateKeyContractAccount(ipp.item.Contract()), getStateFunc); err != nil {
		return errors.Errorf("contract not found, %q; %q", ipp.item.Contract(), err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(ipp.item.Receiver()), getStateFunc); err != nil {
		return errors.Errorf("receiver not found, %q; %q", ipp.item.receiver, err)
	}

	if err := state.CheckNotExistsState(statenft.StateKeyNFT(ipp.item.contract, ipp.idx), getStateFunc); err != nil {
		return errors.Errorf("nft already exists, %q; %q", ipp.idx, err)
	}

	if ipp.item.Creators().Total() != 0 {
		creators := ipp.item.Creators().Signers()
		for _, creator := range creators {
			acc := creator.Account()
			if err := state.CheckExistsState(statecurrency.StateKeyAccount(acc), getStateFunc); err != nil {
				return errors.Errorf("creator not found, %q; %w", acc, err)
			}
			if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(acc), getStateFunc); err != nil {
				return errors.Errorf("contract account cannot be a creator, %q; %w", acc, err)
			}
			if creator.Signed() {
				return errors.Errorf("cannot sign at the same time as minting, %q", acc)
			}
		}
	}

	return nil
}

func (ipp *MintItemProcessor) Process(
	_ context.Context, _ base.Operation, _ base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 1)

	n := types.NewNFT(ipp.idx, true, ipp.item.Receiver(), ipp.item.NFTHash(), ipp.item.URI(), ipp.item.Receiver(), ipp.item.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %q; %w", ipp.idx, err)
	}

	sts[0] = state.NewStateMergeValue(statenft.StateKeyNFT(ipp.item.contract, ipp.idx), statenft.NewNFTStateValue(n))

	if err := ipp.box.Append(n.ID()); err != nil {
		return nil, errors.Errorf("failed to append nft id to nft box, %q; %w", n.ID(), err)
	}

	return sts, nil
}

func (ipp *MintItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = MintItem{}
	ipp.idx = 0
	ipp.box = nil

	mintItemProcessorPool.Put(ipp)

	return nil
}

type MintProcessor struct {
	*base.BaseOperationProcessor
}

func NewMintProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new MintProcessor")

		nopp := mintProcessorPool.Get()
		opp, ok := nopp.(*MintProcessor)
		if !ok {
			return nil, e.Errorf("expected MintProcessor, not %T", nopp)
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

func (opp *MintProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Mint")

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return ctx, nil, e.Errorf("expected MintFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender account not found, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account can not mint nft, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	idxes := map[string]uint64{}
	for _, item := range fact.Items() {
		if _, found := idxes[item.contract.String()]; !found {
			st, err := state.ExistsState(statenft.NFTStateKey(item.contract, statenft.CollectionKey), "key of collection design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection design not found, %q; %w", item.contract, err), nil
			}

			design, err := statenft.StateCollectionValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection design value not found, %q; %w", item.contract, err), nil
			}

			if !design.Active() {
				return nil, base.NewBaseOperationProcessReasonError("deactivated collection, %q", item.contract), nil
			}

			policy, ok := design.Policy().(types.CollectionPolicy)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", types.CollectionPolicy{}, design.Policy()), nil
			}

			whitelist := policy.Whitelist()
			if len(whitelist) == 0 {
				return nil, base.NewBaseOperationProcessReasonError("empty whitelist, %s", item.contract.String()), nil
			}

			st, err = state.ExistsState(stateextension.StateKeyContractAccount(design.Parent()), "key of contract account", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("parent not found, %q; %w", design.Parent(), err), nil
			}

			parent, err := stateextension.StateContractAccountValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("parent value not found, %q; %w", design.Parent(), err), nil
			}

			if !parent.Owner().Equal(fact.Sender()) {
				for i := range whitelist {
					if whitelist[i].Equal(fact.Sender()) {
						break
					}
					if i == len(whitelist)-1 {
						return nil, base.NewBaseOperationProcessReasonError("sender account is neither the owner of contract account nor on the whitelist, %q", fact.Sender()), nil
					}
				}
			}

			if !parent.IsActive() {
				return nil, base.NewBaseOperationProcessReasonError("deactivated parent account, %q", design.Parent()), nil
			}

			st, err = state.ExistsState(statenft.NFTStateKey(item.contract, statenft.LastIDXKey), "key of collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index not found, %q; %w", item.contract, err), nil
			}

			nftID, err := statenft.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index value not found, %q; %w", item.contract, err), nil
			}

			idxes[item.contract.String()] = nftID
		}
	}

	for _, item := range fact.Items() {
		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected MintItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[item.contract.String()]
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess MintItem; %w", err), nil
		}
		idxes[item.contract.String()] += 1

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *MintProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Mint")

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return nil, nil, e.Errorf("expected MintFact, not %T", op.Fact())
	}

	idxes := map[string]uint64{}
	boxes := map[string]*types.NFTBox{}

	for _, item := range fact.items {
		idxKey := statenft.NFTStateKey(item.contract, statenft.LastIDXKey)
		if _, found := idxes[idxKey]; !found {
			st, err := state.ExistsState(idxKey, "key of collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index state not found, %q; %w", item.contract, err), nil
			}

			nftID, err := statenft.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index value not found, %q; %w", item.contract, err), nil
			}

			idxes[idxKey] = nftID
		}

		nftsKey := statenft.NFTStateKey(item.contract, statenft.NFTBoxKey)
		if _, found := boxes[nftsKey]; !found {
			var box types.NFTBox

			switch st, found, err := getStateFunc(nftsKey); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to get nft box state, %q; %w", item.contract, err), nil
			case !found:
				box = types.NewNFTBox(nil)
			default:
				b, err := statenft.StateNFTBoxValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get nft box state value, %q; %w", item.contract, err), nil
				}
				box = b
			}

			boxes[nftsKey] = &box
		}
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*MintItemProcessor, len(fact.Items()))
	for i, item := range fact.Items() {
		idxKey := statenft.NFTStateKey(item.contract, statenft.LastIDXKey)
		nftsKey := statenft.NFTStateKey(item.contract, statenft.NFTBoxKey)
		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected MintItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[idxKey]
		ipc.box = boxes[nftsKey]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process MintItem; %w", err), nil
		}
		sts = append(sts, s...)

		idxes[idxKey] += 1
		ipcs[i] = ipc
	}

	for key, idx := range idxes {
		iv := state.NewStateMergeValue(key, statenft.NewLastNFTIndexStateValue(idx))
		sts = append(sts, iv)
	}

	for key, box := range boxes {
		bv := state.NewStateMergeValue(key, statenft.NewNFTBoxStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	idxes = nil
	boxes = nil

	fitems := fact.Items()
	items := make([]CollectionItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := CalculateCollectionItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, state.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)

	return nil
}
