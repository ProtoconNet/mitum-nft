package collection

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var TransferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferProcessor)
	},
}

func (Transfer) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type TransferProcessor struct {
	cp *currency.CurrencyPool
	Transfer
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewTransferProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Transfer)
		if !ok {
			return nil, errors.Errorf("not Transfer; %T", op)
		}

		opp := TransferProcessorPool.Get().(*TransferProcessor)

		opp.cp = cp
		opp.Transfer = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *TransferProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *TransferProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(TransferFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *TransferProcessor) Close() error {
	opp.cp = nil
	opp.Transfer = Transfer{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	TransferProcessorPool.Put(opp)

	return nil
}
