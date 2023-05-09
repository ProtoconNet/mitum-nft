package collection

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
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
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DelegateItemProcessor struct {
	h      util.Hash
	sender base.Address
	box    *OperatorsBook
	item   DelegateItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if err := checkExistsState(currency.StateKeyAccount(ipp.item.Operator()), getStateFunc); err != nil {
		return err
	}

	if ipp.sender.Equal(ipp.item.Operator()) {
		return errors.Errorf("sender cannot be operator itself, %q", ipp.item.Operator())
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	if ipp.box == nil {
		return nil, errors.Errorf("nft box not found, %q", StateKeyOperators(ipp.item.contract, ipp.item.Collection(), ipp.sender))
	}

	switch ipp.item.Mode() {
	case DelegateAllow:
		if err := ipp.box.Append(ipp.item.Operator()); err != nil {
			return nil, err
		}
	case DelegateCancel:
		if err := ipp.box.Remove(ipp.item.Operator()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("wrong mode for delegate item, %q; \"allow\" | \"cancel\"", ipp.item.Mode())
	}

	ipp.box.Sort(true)

	return nil, nil
}

func (ipp *DelegateItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = DelegateItem{}
	ipp.box = nil

	delegateItemProcessorPool.Put(ipp)

	return nil
}

type DelegateProcessor struct {
	*base.BaseOperationProcessor
}

func NewDelegateProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new DelegateProcessor")

		nopp := delegateProcessorPool.Get()
		opp, ok := nopp.(*DelegateProcessor)
		if !ok {
			return nil, e(nil, "expected DelegateProcessor, not %T", nopp)
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

func (opp *DelegateProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Delegate")

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return ctx, nil, e(nil, "expected DelgateFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot have operators, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, item := range fact.Items() {
		st, err := existsState(NFTStateKey(item.contract, item.Collection(), CollectionKey), "key of design", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("collection design not found, %q: %w", item.Collection(), err), nil
		}

		design, err := StateCollectionValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("collection design value not found, %q: %w", item.Collection(), err), nil
		}

		if !design.Active() {
			return nil, base.NewBaseOperationProcessReasonError("deactivated collection, %q", item.Collection()), nil
		}

		st, err = existsState(extensioncurrency.StateKeyContractAccount(design.Parent()), "key of contract account", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("parent not found, %q: %w", design.Parent(), err), nil
		}

		ca, err := extensioncurrency.StateContractAccountValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %q: %w", design.Parent(), err), nil
		}

		if !ca.IsActive() {
			return nil, base.NewBaseOperationProcessReasonError("deactivated contract account, %q", design.Parent()), nil
		}
	}

	for _, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected DelegateItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess DelegateItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *DelegateProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Delegate")

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return nil, nil, e(nil, "expected DelgateFact, not %T", op.Fact())
	}

	boxes := map[string]*OperatorsBook{}
	for _, item := range fact.Items() {
		ak := StateKeyOperators(item.contract, item.Collection(), fact.Sender())

		var box OperatorsBook
		switch st, found, err := getStateFunc(ak); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError("failed to get state of operators book, %q: %w", ak, err), nil
		case !found:
			box = NewOperatorsBook(item.Collection(), nil)
		default:
			box, err = StateOperatorsBookValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("operators book value not found, %q: %w", ak, err), nil
			}
		}
		boxes[ak] = &box
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*DelegateItemProcessor, len(fact.items))
	for i, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected DelegateItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = boxes[StateKeyOperators(item.contract, item.Collection(), fact.Sender())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process DelegateItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for ak, box := range boxes {
		bv := NewStateMergeValue(ak, NewOperatorsBookStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
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

func (opp *DelegateProcessor) Close() error {
	delegateProcessorPool.Put(opp)

	return nil
}

func (opp *DelegateProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return nil, errors.Errorf("expected DelegateFact, not %T", op.Fact())
	}

	items := make([]CollectionItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateCollectionItemsFee(getStateFunc, items)
}
