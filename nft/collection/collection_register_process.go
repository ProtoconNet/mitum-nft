package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var CollectionRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionRegisterProcessor)
	},
}

func (CollectionRegister) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CollectionRegisterProcessor struct {
	cp *extensioncurrency.CurrencyPool
	CollectionRegister
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewCollectionRegisterProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(CollectionRegister)
		if !ok {
			return nil, errors.Errorf("not CollectionRegister; %T", op)
		}

		opp := CollectionRegisterProcessorPool.Get().(*CollectionRegisterProcessor)

		opp.cp = cp
		opp.CollectionRegister = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *CollectionRegisterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *CollectionRegisterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(CollectionRegisterFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *CollectionRegisterProcessor) Close() error {
	opp.cp = nil
	opp.CollectionRegister = CollectionRegister{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	CollectionRegisterProcessorPool.Put(opp)

	return nil
}
