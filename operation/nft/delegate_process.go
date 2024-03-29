package nft

import (
	"context"
	"sync"

	statenft "github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var delegateItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateItemProcessor)
	},
}

var delegateProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateProcessor)
	},
}

func (Delegate) Process(
	ctx context.Context, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DelegateItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	box    *types.OperatorsBook
	item   DelegateItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(ipp.item.Delegatee()), getStateFunc); err != nil {
		return err
	}

	if ipp.sender.Equal(ipp.item.Delegatee()) {
		return errors.Errorf("sender account cannot delegate to itself, %q", ipp.item.Delegatee())
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	if ipp.box == nil {
		return nil, errors.Errorf("nft box not found, %q", statenft.StateKeyOperators(ipp.item.Contract(), ipp.sender))
	}

	switch ipp.item.Mode() {
	case DelegateAllow:
		if err := ipp.box.Append(ipp.item.Delegatee()); err != nil {
			return nil, err
		}
	case DelegateCancel:
		if err := ipp.box.Remove(ipp.item.Delegatee()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("wrong mode for delegate item, %q; \"allow\" | \"cancel\"", ipp.item.Mode())
	}

	ipp.box.Sort(true)

	return nil, nil
}

func (ipp *DelegateItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = DelegateItem{}
	ipp.box = nil

	delegateItemProcessorPool.Put(ipp)

	return
}

type DelegateProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewDelegateProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new DelegateProcessor")

		nopp := delegateProcessorPool.Get()
		opp, ok := nopp.(*DelegateProcessor)
		if !ok {
			return nil, e.Errorf("expected DelegateProcessor, not %T", nopp)
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

func (opp *DelegateProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Delegate")

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return ctx, nil, e.Errorf("expected %T, not %T", DelegateFact{}, op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender account not found, %q; %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("contract account cannot have operators, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	for _, item := range fact.Items() {
		err := state.CheckExistsState(stateextension.StateKeyContractAccount(item.Contract()), getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("contract account not found, %q; %w", item.Contract(), err), nil
		}

		st, err := state.ExistsState(statenft.NFTStateKey(item.contract, statenft.CollectionKey), "key of design", getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("collection design not found, %q; %w", item.Contract(), err), nil
		}

		design, err := statenft.StateCollectionValue(st)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("collection design value not found, %q; %w", item.Contract(), err), nil
		}

		if !design.Active() {
			return nil, mitumbase.NewBaseOperationProcessReasonError("deactivated collection, %q", item.Contract()), nil
		}
	}

	for _, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected DelegateItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("fail to preprocess DelegateItem; %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *DelegateProcessor) Process(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Delegate")

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return nil, nil, e.Errorf("expected DelgateFact, not %T", op.Fact())
	}

	boxes := map[string]*types.OperatorsBook{}
	for _, item := range fact.Items() {
		ak := statenft.StateKeyOperators(item.contract, fact.Sender())

		var operators types.OperatorsBook
		switch st, found, err := getStateFunc(ak); {
		case err != nil:
			return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get state of operators book, %q; %w", ak, err), nil
		case !found:
			operators = types.NewOperatorsBook(nil)
		default:
			o, err := statenft.StateOperatorsBookValue(st)
			if err != nil {
				return nil, mitumbase.NewBaseOperationProcessReasonError("operators book value not found, %q; %w", ak, err), nil
			} else {
				operators = *o
			}
		}
		boxes[ak] = &operators
	}

	var sts []mitumbase.StateMergeValue // nolint:prealloc

	ipcs := make([]*DelegateItemProcessor, len(fact.items))
	for i, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e.Errorf("expected %T, not %T", &DelegateItemProcessor{}, ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = boxes[statenft.StateKeyOperators(item.contract, fact.Sender())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError("failed to process DelegateItem; %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for ak, box := range boxes {
		bv := state.NewStateMergeValue(ak, statenft.NewOperatorsBookStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	items := make([]CollectionItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
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

func (opp *DelegateProcessor) Close() error {
	delegateProcessorPool.Put(opp)

	return nil
}
