package collection

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

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

type DelegateProcessor struct {
	cp *currency.CurrencyPool
	Delegate
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewDelegateProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Delegate)
		if !ok {
			return nil, errors.Errorf("not Delegate; %T", op)
		}

		opp := DelegateProcessorPool.Get().(*DelegateProcessor)

		opp.cp = cp
		opp.Delegate = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *DelegateProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *DelegateProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(DelegateFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *DelegateProcessor) Close() error {
	opp.cp = nil
	opp.Delegate = Delegate{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	DelegateProcessorPool.Put(opp)

	return nil
}
