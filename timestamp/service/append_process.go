package service

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var appendProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AppendProcessor)
	},
}

func (Append) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AppendProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc GetLastBlockFunc
}

func NewAppendProcessor(getLastBlockFunc GetLastBlockFunc) extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new AppendProcessor")

		nopp := appendProcessorPool.Get()
		opp, ok := nopp.(*AppendProcessor)
		if !ok {
			return nil, e(nil, "expected AppendProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b
		opp.getLastBlockFunc = getLastBlockFunc

		return opp, nil
	}
}

func (opp *AppendProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Mint")

	fact, ok := op.Fact().(AppendFact)
	if !ok {
		return ctx, nil, e(nil, "expected AppendFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot Append timestamp, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	_, err := existsState(StateKeyServiceDesign(fact.target, fact.service), "key of service design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design not found, %q: %w", fact.service, err), nil
	}

	k := StateKeyTimeStampLastIndex(fact.target, fact.service, fact.projectID)
	switch _, _, err := getStateFunc(k); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("getting timestamp item lastindex failed, %q: %w", fact.service, err), nil
	}

	_, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed, %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	return ctx, nil, nil
}

func (opp *AppendProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Append")

	fact, ok := op.Fact().(AppendFact)
	if !ok {
		return nil, nil, e(nil, "expected AppendFact, not %T", op.Fact())
	}

	st, err := existsState(StateKeyServiceDesign(fact.target, fact.service), "key of service design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design not found, %q: %w", fact.service, err), nil
	}

	design, err := StateServiceDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design value not found, %q: %w", fact.service, err), nil
	}

	design.AddProject(fact.projectID)

	var idx uint64
	k := StateKeyTimeStampLastIndex(fact.target, fact.service, fact.projectID)
	switch st, found, err := getStateFunc(k); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("getting timestamp item lastindex failed, %q: %w", fact.service, err), nil
	case found:
		i, err := StateTimeStampLastIndexValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("getting timestamp item lastindex value failed, %q: %w", fact.service, err), nil
		}
		idx = i + 1
	case !found:
		idx = 0
		st = base.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}

	blockmap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed, %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	tsItem := timestamp.NewTimeStampItem(fact.projectID, fact.requestTimeStamp, uint64(blockmap.Manifest().ProposedAt().Unix()), idx, fact.data)
	if err := tsItem.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid timestamp, %w", err), nil
	}

	sts := make([]base.StateMergeValue, 2) // nolint:prealloc
	sts[0] = NewStateMergeValue(StateKeyTimeStampItem(fact.target, fact.service, fact.projectID, idx), NewTimeStampItemStateValue(tsItem))
	sts[1] = NewStateMergeValue(StateKeyTimeStampLastIndex(fact.target, fact.service, fact.projectID), NewTimeStampLastIndexStateValue(fact.projectID, idx))

	return sts, nil, nil
}

func (opp *AppendProcessor) Close() error {
	appendProcessorPool.Put(opp)

	return nil
}
