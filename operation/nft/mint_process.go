package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
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
	ns     map[string]base.StateMergeValue
}

func (ipp *MintItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess MintItemProcessor")

	if err := currencystate.CheckExistsState(
		statecurrency.DesignStateKey(ipp.item.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", ipp.item.Currency())))
	}

	if _, _, _, cErr := currencystate.ExistsCAccount(
		ipp.item.receiver, "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: receiver %v is contract account", cErr, ipp.item.Receiver())))
	}

	if found, _ := currencystate.CheckNotExistsState(
		statenft.StateKeyNFT(ipp.item.Contract(), ipp.idx), getStateFunc); found {
		return e.Wrap(
			common.ErrStateE.Wrap(
				errors.Errorf("nft idx %v already exists in contract account %v", ipp.idx, ipp.item.Contract())))
	}

	creators := ipp.item.Creators().Signers()
	for _, creator := range creators {
		acc := creator.Address()
		if _, _, _, cErr := currencystate.ExistsCAccount(
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

	smv, err := currencystate.CreateNotExistAccount(ipp.item.Receiver(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	creators := ipp.item.Creators().Signers()
	for _, creator := range creators {
		smv, err := currencystate.CreateNotExistAccount(creator.Address(), getStateFunc)
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

	sts = append(sts, currencystate.NewStateMergeValue(statenft.StateKeyNFT(ipp.item.Contract(), ipp.idx), statenft.NewNFTStateValue(n)))

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

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender %v is contract account", fact.Sender(), cErr)), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	idxes := map[string]uint64{}
	for _, item := range fact.Items() {
		if _, found := idxes[item.contract.String()]; !found {
			_, _, aErr, cErr := currencystate.ExistsCAccount(
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

			st, err := currencystate.ExistsState(
				statenft.NFTStateKey(item.contract, statenft.CollectionKey), "design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf("nft collection, %v: %v", item.Contract(), err)), nil
			}

			design, err := statenft.StateCollectionValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf("nft collection, %v: %v", item.Contract(), err)), nil
			}

			if !design.Active() {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf(
						"collection in contract account %v has already been deactivated", item.Contract())), nil

			}

			policy, ok := design.Policy().(types.CollectionPolicy)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMTypeMismatch).
						Errorf("expected %T, not %T", types.CollectionPolicy{}, design.Policy())), nil
			}

			whitelist := policy.Whitelist()
			_, cSt, aErr, cErr := currencystate.ExistsCAccount(
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

			ca, err := stateextension.LoadCAStateValue(cSt)
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

			if !ca.IsActive() {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMServiceNF).Errorf("nft collection, %v", item.Contract())), nil
			}

			st, err = currencystate.ExistsState(statenft.NFTStateKey(item.contract, statenft.LastIDXKey), "collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.
						Wrap(common.ErrMStateNF).Errorf("collection last index, %v: %v", item.contract, err)), nil
			}

			nftID, err := statenft.StateLastNFTIndexValue(st)
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

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return nil, nil, e.Errorf("expected MintFact, not %T", op.Fact())
	}

	idxes := map[string]uint64{}
	//boxes := map[string]*types.NFTBox{}

	for _, item := range fact.items {
		idxKey := statenft.NFTStateKey(item.contract, statenft.LastIDXKey)
		if _, found := idxes[idxKey]; !found {
			st, err := currencystate.ExistsState(idxKey, "collection index", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("collection last index state not found, %v: %w", item.contract, err), nil
			}

			nftID, err := statenft.StateLastNFTIndexValue(st)
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
		idxKey := statenft.NFTStateKey(item.contract, statenft.LastIDXKey)
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
		iv := currencystate.NewStateMergeValue(key, statenft.NewLastNFTIndexStateValue(idx))
		sts = append(sts, iv)
	}

	for _, ns := range nsts {
		sts = append(sts, ns)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	idxes = nil
	//boxes = nil

	items := make([]CollectionItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiverBalSts, required, err := CalculateCollectionItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.Sender(), required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
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
				func(height base.Height, st base.State) base.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiverBalSts[cid].Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", statecurrency.BalanceStateValue{}, feeReceiverBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiverBalSts[cid].Key(),
					statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height base.Height, st base.State) base.StateValueMerger {
						return statecurrency.NewBalanceStateValueMerger(height, feeReceiverBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)

	return nil
}
