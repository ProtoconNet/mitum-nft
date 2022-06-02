package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var ApproveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type ApproveProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Approve
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewApproveProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Approve)
		if !ok {
			return nil, errors.Errorf("not Approve; %T", op)
		}

		opp := ApproveProcessorPool.Get().(*ApproveProcessor)

		opp.cp = cp
		opp.Approve = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *ApproveProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *ApproveProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(ApproveFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *ApproveProcessor) Close() error {
	opp.cp = nil
	opp.Approve = Approve{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	ApproveProcessorPool.Put(opp)

	return nil
}
