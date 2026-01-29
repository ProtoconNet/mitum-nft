package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	cestate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
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
	ns     map[string]base.StateMergeValue
}

func (ipp *MintItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess MintItemProcessor")

	if err := ipp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		ipp.item.receiver, "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: receiver %v is contract account", cErr, ipp.item.Receiver())))
	}

	if found, _ := cstate.CheckNotExistsState(
		state.StateKeyNFT(ipp.item.Contract(), ipp.idx), getStateFunc); found {
		return e.Wrap(
			common.ErrStateE.Wrap(
				errors.Errorf("nft idx %v already exists in contract account %v", ipp.idx, ipp.item.Contract())))
	}

	creators := ipp.item.Creators().Signers()
	for _, creator := range creators {
		acc := creator.Address()
		if _, _, _, cErr := cstate.ExistsCAccount(
			acc, "creator", true, false, getStateFunc); cErr != nil {
			return e.Wrap(common.ErrCAccountNA.Wrap(
				errors.Errorf("%v: creator %v is contract account", cErr, acc)))
		}
		if creator.Signed() {
			return e.Wrap(errors.Errorf("creator %v must be unsigned at the time of minting", acc))
		}
	}

	return nil
}

func (ipp *MintItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(ipp.item.Receiver(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	creators := ipp.item.Creators().Signers()
	for _, creator := range creators {
		smv, err := cstate.CreateNotExistAccount(creator.Address(), getStateFunc)
		if err != nil {
			return nil, err
		} else if smv != nil {
			sts = append(sts, smv)
		}
	}

	n := types.NewNFT(ipp.idx, true, ipp.item.Receiver(), ipp.item.NFTHash(), ipp.item.URI(), ipp.item.Receiver(), ipp.item.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, errors.Errorf("invalid nft, %v: %v", ipp.idx, err)
	}

	sts = append(sts, cstate.NewStateMergeValue(state.StateKeyNFT(ipp.item.Contract(), ipp.idx), state.NewNFTStateValue(n)))

	st, _ := cstate.ExistsState(state.NFTStateKey(ipp.item.contract, state.CollectionKey), "design", getStateFunc)
	design, _ := state.StateCollectionValue(st)
	de := types.NewDesign(design.Contract(), design.Creator(), design.Active(), design.Count()+1, design.Policy())
	sts = append(
		sts,
		cstate.NewStateMergeValue(
			state.NFTStateKey(ipp.item.contract, state.CollectionKey),
			state.NewCollectionStateValue(de)),
	)

	return sts, nil
}

func (ipp *MintItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = MintItem{}
	ipp.idx = 0
	//ipp.box = nil
	ipp.ns = nil

	mintItemProcessorPool.Put(ipp)

	return
}

type MintProcessor struct {
	*base.BaseOperationProcessor
}

func NewMintProcessor() ctypes.GetNewProcessor {
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
	fact, ok := op.Fact().(MintFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", MintFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	idxes := map[string]uint64{}
	for _, item := range fact.Items() {
		if _, found := idxes[item.contract.String()]; !found {
			st, err := cstate.ExistsState(
				state.NFTStateKey(item.contract, state.CollectionKey), "design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf("nft service state for contract account, %v: %v", item.Contract(), err)), nil
			}

			design, err := state.StateCollectionValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf("nft service state for contract account, %v: %v", item.Contract(), err)), nil
			}

			if !design.Active() {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf(
						"nft service in contract account %v has already been deactivated", item.Contract())), nil

			}

			policy, ok := design.Policy().(types.CollectionPolicy)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMTypeMismatch).
						Errorf("expected %T, not %T", types.CollectionPolicy{}, design.Policy())), nil
			}

			whitelist := policy.Whitelist()
			_, cSt, aErr, cErr := cstate.ExistsCAccount(
				item.Contract(), "contract", true, true, getStateFunc)
			if aErr != nil {
				return ctx, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Errorf("%v", aErr)), nil
			} else if cErr != nil {
				return ctx, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Errorf("%v", cErr)), nil
			}

			ca, err := cestate.LoadCAStateValue(cSt)
			if err != nil {
				return ctx, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Errorf("%v", err)), nil
			}
			if !ca.Owner().Equal(fact.Sender()) {
				for i := range whitelist {
					if whitelist[i].Equal(fact.Sender()) {
						break
					}
					if i == len(whitelist)-1 {
						return ctx, base.NewBaseOperationProcessReasonError(
							common.ErrMPreProcess.Wrap(common.ErrMAccountNAth).
								Errorf(
									"sender %v is neither the owner nor in the minter whitelist of contract account %v",
									fact.Sender(), item.Contract())), nil
					}
				}
			}

			st, err = cstate.ExistsState(state.NFTStateKey(item.contract, state.LastIDXKey), "collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMStateNF).Errorf("collection last index, %v: %v", item.contract, err)), nil
			}

			nftID, err := state.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMStateInvalid).Errorf("collection last index, %v: %v", item.contract, err)), nil
			}

			idxes[item.contract.String()] = nftID
		}
	}

	for _, item := range fact.Items() {
		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected MintItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[item.contract.String()]
		//ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
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

	fact, _ := op.Fact().(MintFact)
	idxes := map[string]uint64{}

	for _, item := range fact.items {
		idxKey := state.NFTStateKey(item.contract, state.LastIDXKey)
		if _, found := idxes[idxKey]; !found {
			st, err := cstate.ExistsState(idxKey, "collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index state not found, %v: %w", item.contract, err), nil
			}

			nftID, err := state.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index value not found, %v: %w", item.contract, err), nil
			}

			idxes[idxKey] = nftID
		}
	}

	var sts []base.StateMergeValue // nolint:prealloc
	nsts := map[string]base.StateMergeValue{}

	ipcs := make([]*MintItemProcessor, len(fact.Items()))
	for i, item := range fact.Items() {
		idxKey := state.NFTStateKey(item.contract, state.LastIDXKey)
		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected MintItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[idxKey]
		ipc.ns = nsts

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process MintItem; %w", err), nil
		}
		sts = append(sts, s...)

		idxes[idxKey] += 1
		ipcs[i] = ipc
	}

	for key, idx := range idxes {
		iv := cstate.NewStateMergeValue(key, state.NewLastNFTIndexStateValue(idx))
		sts = append(sts, iv)
	}

	for _, ns := range nsts {
		sts = append(sts, ns)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	idxes = nil

	return sts, nil, nil
}

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)

	return nil
}
