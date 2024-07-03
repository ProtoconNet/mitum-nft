package nft

import (
	"context"
	"sync"

	statenft "github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
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

func (ApproveAll) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DelegateItemProcessor struct {
	h      util.Hash
	sender mitumbase.Address
	box    *types.AllApprovedBook
	item   ApproveAllItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	_ context.Context, _ mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) error {
	e := util.StringError("preprocess DelegateItemProcessor")

	if err := currencystate.CheckExistsState(
		statecurrency.DesignStateKey(ipp.item.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", ipp.item.Currency())))
	}

	if _, _, _, cErr := currencystate.ExistsCAccount(
		ipp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(
			errors.Errorf("%v: approved %v is contract account", cErr, ipp.item.Approved())))
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	_ context.Context, _ mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, error) {
	var sts []mitumbase.StateMergeValue

	smv, err := currencystate.CreateNotExistAccount(ipp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, err
	} else if smv != nil {
		sts = append(sts, smv)
	}

	if ipp.box == nil {
		return nil, errors.Errorf(
			"nft box not found, %v", statenft.StateKeyOperators(ipp.item.Contract(), ipp.sender))
	}

	switch ipp.item.Mode() {
	case ApproveAllAllow:
		if err := ipp.box.Append(ipp.item.Approved()); err != nil {
			return nil, err
		}
	case ApproveAllCancel:
		if err := ipp.box.Remove(ipp.item.Approved()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("wrong mode for delegate item, %v: \"allow\" | \"cancel\"", ipp.item.Mode())
	}

	ipp.box.Sort(true)

	return sts, nil
}

func (ipp *DelegateItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveAllItem{}
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
	fact, ok := op.Fact().(ApproveAllFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ApproveAllFact{}, op.Fact())), nil
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
				Errorf("%v: sender %v is contract", fact.Sender(), cErr)), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, item := range fact.Items() {
		_, _, aErr, cErr := currencystate.ExistsCAccount(
			item.Contract(), "contract", true, true, getStateFunc)
		if aErr != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", aErr)), nil
		} else if cErr != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", cErr)), nil
		}

		st, err := currencystate.ExistsState(
			statenft.NFTStateKey(item.Contract(), statenft.CollectionKey), "design", getStateFunc)
		if err != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("nft collection, %v: %v", item.Contract(), err)), nil
		}

		design, err := statenft.StateCollectionValue(st)
		if err != nil {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("nft collection, %v: %v", item.Contract(), err)), nil
		}

		if !design.Active() {
			return ctx, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
					Errorf("collection in contract account %v has already been deactivated",
						item.Contract())), nil
		}
	}

	for _, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected DelegateItemProcessor, not %T", ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
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

	fact, ok := op.Fact().(ApproveAllFact)
	if !ok {
		return nil, nil, e.Errorf("expected DelgateFact, not %T", op.Fact())
	}

	boxes := map[string]*types.AllApprovedBook{}
	for _, item := range fact.Items() {
		ak := statenft.StateKeyOperators(item.contract, fact.Sender())

		var operators types.AllApprovedBook
		switch st, found, err := getStateFunc(ak); {
		case err != nil:
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				"failed to get state of operators book, %v: %w", ak, err), nil
		case !found:
			operators = types.NewAllApprovedBook(nil)
		default:
			o, err := statenft.StateOperatorsBookValue(st)
			if err != nil {
				return nil, mitumbase.NewBaseOperationProcessReasonError(
					"operators book value not found, %v: %w", ak, err), nil
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
