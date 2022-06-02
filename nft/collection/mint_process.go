package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MintProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintProcessor)
	},
}

func (Mint) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type MintProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Mint
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewMintProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Mint)
		if !ok {
			return nil, errors.Errorf("not Mint; %T", op)
		}

		opp := MintProcessorPool.Get().(*MintProcessor)

		opp.cp = cp
		opp.Mint = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *MintProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *MintProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(MintFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *MintProcessor) Close() error {
	opp.cp = nil
	opp.Mint = Mint{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	MintProcessorPool.Put(opp)

	return nil
}
