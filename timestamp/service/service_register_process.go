package service

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var serviceRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ServiceRegisterProcessor)
	},
}

func (ServiceRegister) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ServiceRegisterProcessor struct {
	*base.BaseOperationProcessor
}

func NewServiceRegisterProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new ServiceRegisterProcessor")

		nopp := serviceRegisterProcessorPool.Get()
		opp, ok := nopp.(*ServiceRegisterProcessor)
		if !ok {
			return nil, errors.Errorf("expected servicesRegisterProcessor, not %T", nopp)
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

func (opp *ServiceRegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess serviceRegister")

	fact, ok := op.Fact().(ServiceRegisterFact)
	if !ok {
		return ctx, nil, e(nil, "expected ServiceRegisterFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender address is contract account, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	st, err := existsState(extensioncurrency.StateKeyContractAccount(fact.Target()), "key of contract account", getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("target contract account not found, %q: %w", fact.Target(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q: %w", fact.Target(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return ctx, base.NewBaseOperationProcessReasonError("sender is not owner of contract account, %q, %q", fact.Sender(), ca.Owner()), nil
	}

	if !ca.IsActive() {
		return ctx, base.NewBaseOperationProcessReasonError("deactivated contract account, %q", fact.Target()), nil
	}

	if err := checkNotExistsState(StateKeyServiceDesign(fact.Target(), fact.Service()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("service design already exists, %q: %w", fact.Service(), err), nil
	}

	return ctx, nil, nil
}

func (opp *ServiceRegisterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process ServiceRegister")

	fact, ok := op.Fact().(ServiceRegisterFact)
	if !ok {
		return nil, nil, e(nil, "expected ServiceRegisterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 2)
	pids := []string{}

	design := timestamp.NewDesign(fact.Service(), pids...)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q: %w", fact.Service(), err), nil
	}

	sts[0] = NewStateMergeValue(
		StateKeyServiceDesign(fact.target, design.Service()),
		NewServiceDesignStateValue(design),
	)

	currencyPolicy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err := existsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := currency.NewBalanceStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = currency.NewBalanceStateMergeValue(
		sb.Key(),
		currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *ServiceRegisterProcessor) Close() error {
	serviceRegisterProcessorPool.Put(opp)

	return nil
}
