package nft

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	item   ApproveItem
}

func (ipp *ApproveItemProcessor) PreProcess(
	_ context.Context, _ mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	e := util.StringError("preprocess ApproveItemProcessor")

	if err := currencystate.CheckExistsState(
		statecurrency.StateKeyCurrencyDesign(ipp.item.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id, %v", ipp.item.Currency())))
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(
		ipp.item.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(errors.Errorf("%v", cErr))
	}

	st, err := state.ExistsState(
		statenft.NFTStateKey(ipp.item.Contract(), statenft.CollectionKey), "design", getStateFunc)
	if err != nil {
		return e.Wrap(
			common.ErrServiceNF.Wrap(errors.Errorf("nft collection, %v: %v", ipp.item.Contract(), err)))
	}

	design, err := statenft.StateCollectionValue(st)
	if err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("nft collection, %v: %v", ipp.item.Contract(), err)))
	}

	if !design.Active() {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("deactivated collection, %v", ipp.item.Contract())))
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		ipp.item.Approved(), "approved", true, false, getStateFunc); aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(errors.Errorf("%v: contract account, %v cannot be approved", cErr, ipp.item.Approved()))
	}

	st, err = state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), ipp.item.idx), "nft", getStateFunc)
	if err != nil {
		return common.ErrStateNF.Wrap(errors.Errorf("nft, %v: %v", ipp.item.idx, err))
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return common.ErrStateValInvalid.Wrap(errors.Errorf("nft, %v: %v", ipp.item.idx, err))
	}

	if !nv.Active() {
		return common.ErrStateValInvalid.Wrap(errors.Errorf("burned nft, %v", ipp.item.idx))
	}

	if ipp.item.Approved().Equal(nv.Approved()) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("already approved, %v", ipp.item.Approved()))
	}

	if !nv.Owner().Equal(ipp.sender) {
		if err := state.CheckExistsState(statecurrency.StateKeyAccount(nv.Owner()), getStateFunc); err != nil {
			return common.ErrAccountNF.Wrap(errors.Errorf("nft owner, %v: %v", nv.Owner(), err))
		}

		st, err = state.ExistsState(statenft.StateKeyOperators(ipp.item.Contract(), ipp.sender), "operators", getStateFunc)
		if err != nil {
			return common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v neither nft owner nor operator: %v", ipp.sender, err))
		}

		operators, err := statenft.StateOperatorsBookValue(st)
		if err != nil {
			return common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v neither nft owner nor operator: %v", ipp.sender, err))
		}

		if !operators.Exists(ipp.sender) {
			return common.ErrAccountNAth.Wrap(errors.Errorf("sender, %v: neither nft owner nor operator", ipp.sender))
		}
	}

	return nil
}

func (ipp *ApproveItemProcessor) Process(
	_ context.Context, _ mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := state.ExistsState(statenft.StateKeyNFT(ipp.item.Contract(), nid), "nft", getStateFunc)
	if err != nil {
		return nil, util.ErrNotFound.Errorf("nft, %v: %v", nid, err)
	}

	nv, err := statenft.StateNFTValue(st)
	if err != nil {
		return nil, util.ErrNotFound.Errorf("nft value, %v: %v", nid, err)
	}

	n := types.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), ipp.item.Approved(), nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []mitumbase.StateMergeValue{state.NewStateMergeValue(st.Key(), statenft.NewNFTStateValue(n))}

	return sts, nil
}

func (ipp *ApproveItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveItem{}

	approveItemProcessorPool.Put(ipp)

	return
}

type ApproveProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewApproveProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new ApproveProcessor")

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
		if !ok {
			return nil, e.Errorf("expected ApproveProcessor, not %T", nopp)
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ApproveFact{}, op.Fact())), nil
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

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected ApproveItemProcessor, not %T", ip)), nil
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

func (opp *ApproveProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Approve")

	fact, _ := op.Fact().(ApproveFact)

	var sts []mitumbase.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected ApproveItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("process ApproveItem: %w", err), nil
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
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to calculate fee: %v", err), nil
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

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)

	return nil
}

func CalculateCollectionItemsFee(getStateFunc mitumbase.GetStateFunc, items []CollectionItem) (
	map[currencytypes.CurrencyID]mitumbase.State, map[currencytypes.CurrencyID][2]common.Big, error) {
	feeReceiveSts := map[currencytypes.CurrencyID]mitumbase.State{}
	required := map[currencytypes.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := state.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

		if policy.Feeer().Receiver() == nil {
			continue
		}

		if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(policy.Feeer().Receiver()), getStateFunc); err != nil {
			return nil, nil, err
		} else if st, found, err := getStateFunc(statecurrency.StateKeyBalance(policy.Feeer().Receiver(), item.Currency())); err != nil {
			return nil, nil, err
		} else if !found {
			return nil, nil, errors.Errorf("feeer receiver account not found, %s", policy.Feeer().Receiver())
		} else {
			feeReceiveSts[item.Currency()] = st
		}

	}

	return feeReceiveSts, required, nil

}
