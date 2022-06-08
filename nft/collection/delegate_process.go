package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var DelegateItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateItemProcessor)
	},
}

var DelegateProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateProcessor)
	},
}

func (Delegate) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type DelegateItemProcessor struct {
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	box    *AgentBox
	sender base.Address
	item   DelegateItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if err := checkExistsState(currency.StateKeyAccount(ipp.item.Agent()), getState); err != nil {
		return err
	}

	if ipp.sender.Equal(ipp.item.Agent()) {
		return errors.Errorf("sender cannot be agent itself; %q", ipp.item.Agent().String())
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	switch ipp.item.Mode() {
	case DelegateAllow:
		if err := ipp.box.Append(ipp.item.Agent()); err != nil {
			return nil, err
		}
	case DelegateCancel:
		if err := ipp.box.Remove(ipp.item.Agent()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("wrong mode for delegate item; mode must be [\"allow\": delegate | \"cancel\": cancel delegation]")
	}

	return nil, nil
}

func (ipp *DelegateItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.sender = nil
	ipp.item = DelegateItem{}
	ipp.box = nil

	DelegateItemProcessorPool.Put(ipp)

	return nil
}

type DelegateProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Delegate
	box          AgentBox
	boxState     state.State
	amountStates map[currency.CurrencyID]currency.AmountState
	ipps         []*DelegateItemProcessor
	required     map[currency.CurrencyID][2]currency.Big
}

func NewDelegateProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Delegate)
		if !ok {
			return nil, operation.NewBaseReasonError("not Delegate; %T", op)
		}

		opp := DelegateProcessorPool.Get().(*DelegateProcessor)

		opp.cp = cp
		opp.Delegate = i
		opp.box = AgentBox{}
		opp.boxState = nil
		opp.amountStates = nil
		opp.ipps = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *DelegateProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(DelegateFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot have agents; %q", fact.Sender())
	}

	switch st, found, err := getState(StateKeyAgents(fact.Sender())); {
	case err != nil:
		return nil, operation.NewBaseReasonError(err.Error())
	case !found:
		opp.box = NewAgentBox(nil)
		opp.boxState = st
	default:
		box, err := StateAgentsValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}
		opp.box = box
		opp.boxState = st
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else if sts, err := CheckSenderEnoughBalance(fact.Sender(), required, getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.required = required
		opp.amountStates = sts
	}

	ipps := make([]*DelegateItemProcessor, len(fact.items))
	for i := range fact.items {

		c := DelegateItemProcessorPool.Get().(*DelegateItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.sender = fact.Sender()
		c.item = fact.items[i]
		c.box = &opp.box

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}

		ipps[i] = c
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	opp.ipps = ipps

	return opp, nil
}

func (opp *DelegateProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(DelegateFact)

	var states []state.State

	for i := range opp.ipps {
		if s, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process delegate item; %w", err)
		} else {
			states = append(states, s...)
		}
	}
	opp.box.Sort(true)

	if st, err := SetStateAgentsValue(opp.boxState, opp.box); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	for k := range opp.required {
		rq := opp.required[k]
		states = append(states, opp.amountStates[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), states...)
}

func (opp *DelegateProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Delegate = Delegate{}
	opp.box = AgentBox{}
	opp.boxState = nil
	opp.amountStates = nil
	opp.ipps = nil
	opp.required = nil

	DelegateProcessorPool.Put(opp)

	return nil
}

func (opp *DelegateProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(DelegateFact)

	items := make([]DelegateItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateDelegateItemsFee(opp.cp, items)
}

func CalculateDelegateItemsFee(cp *extensioncurrency.CurrencyPool, items []DelegateItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for i := range items {
		it := items[i]

		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}

		if cp == nil {
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}
			continue
		}

		feeer, found := cp.Feeer(it.Currency())
		if !found {
			return nil, errors.Errorf("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[it.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}

func CheckSenderEnoughBalance(
	holder base.Address,
	required map[currency.CurrencyID][2]currency.Big,
	getState func(key string) (state.State, bool, error),
) (map[currency.CurrencyID]currency.AmountState, error) {
	sb := map[currency.CurrencyID]currency.AmountState{}

	for cid := range required {
		rq := required[cid]

		st, err := existsState(currency.StateKeyBalance(holder, cid), "currency of holder", getState)
		if err != nil {
			return nil, err
		}

		am, err := currency.StateBalanceValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError("insufficient balance of sender: %w", err)
		}

		if am.Big().Compare(rq[0]) < 0 {
			return nil, operation.NewBaseReasonError(
				"insufficient balance of sender, %s; %d !> %d", holder.String(), am.Big(), rq[0])
		} else {
			sb[cid] = currency.NewAmountState(st, cid)
		}
	}

	return sb, nil
}
