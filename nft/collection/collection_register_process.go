package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/operation"
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
	idxState    state.State
	DesignState state.State
	design      nft.Design
	amountState currency.AmountState
	fee         currency.Big
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
		opp.idxState = nil
		opp.DesignState = nil
		opp.design = nft.Design{}
		opp.amountState = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *CollectionRegisterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact, ok := opp.Fact().(CollectionRegisterFact)
	if !ok {
		return nil, operation.NewBaseReasonError("not CollectionRegisterFact; %T", opp.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot register a collection; %q", fact.Sender())
	}

	if st, err := existsState(extensioncurrency.StateKeyContractAccount(fact.Form().Target()), "contract account", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if ca, err := extensioncurrency.StateContractAccountValue(st); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if !ca.Owner().Equal(fact.Sender()) {
		return nil, operation.NewBaseReasonError("not owner of contract account; %q", fact.Form().Target())
	} else if !ca.IsActive() {
		return nil, operation.NewBaseReasonError("deactivated contract account; %q", fact.Form().Target())
	}

	if st, err := notExistsState(StateKeyCollection(fact.Form().Symbol()), "design", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.DesignState = st
	}

	if st, err := notExistsState(StateKeyCollectionLastIDX(fact.Form().Symbol()), "collection idx", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.idxState = st
	}

	whites := fact.Form().Whites()
	for i := range whites {
		if err := checkExistsState(currency.StateKeyAccount(whites[i]), getState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		} else if err = checkNotExistsState(extensioncurrency.StateKeyContractAccount(whites[i]), getState); err != nil {
			return nil, operation.NewBaseReasonError("contract account cannot be whitelisted; %q", whites[i])
		}
	}

	policy := NewCollectionPolicy(fact.Form().Name(), fact.Form().Royalty(), fact.Form().Uri(), whites)
	if err := policy.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	design := nft.NewDesign(fact.Form().Target(), fact.Sender(), fact.Form().Symbol(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}
	opp.design = design

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	if st, err := existsState(
		currency.StateKeyBalance(fact.Sender(), fact.Currency()), "balance of sender", getState); err != nil {
		return nil, err
	} else {
		opp.amountState = currency.NewAmountState(st, fact.Currency())
	}

	feeer, found := opp.cp.Feeer(fact.Currency())
	if !found {
		return nil, operation.NewBaseReasonError("currency not found; %q", fact.Currency())
	}

	fee, err := feeer.Fee(currency.ZeroBig)
	if err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	}
	switch b, err := currency.StateBalanceValue(opp.amountState); {
	case err != nil:
		return nil, operation.NewBaseReasonErrorFromError(err)
	case b.Big().Compare(fee) < 0:
		return nil, operation.NewBaseReasonError("insufficient balance with fee")
	default:
		opp.fee = fee
	}

	return opp, nil
}

func (opp *CollectionRegisterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact, ok := opp.Fact().(CollectionRegisterFact)
	if !ok {
		return operation.NewBaseReasonError("not CollectionRegisterFact; %T", opp.Fact())
	}

	var states []state.State

	if st, err := SetStateCollectionLastIDXValue(opp.idxState, 0); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	if st, err := SetStateCollectionValue(opp.DesignState, opp.design); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	opp.amountState = opp.amountState.Sub(opp.fee).AddFee(opp.fee)
	states = append(states, opp.amountState)

	return setState(fact.Hash(), states...)
}

func (opp *CollectionRegisterProcessor) Close() error {
	opp.cp = nil
	opp.CollectionRegister = CollectionRegister{}
	opp.idxState = nil
	opp.DesignState = nil
	opp.design = nft.Design{}
	opp.amountState = currency.AmountState{}
	opp.fee = currency.ZeroBig

	CollectionRegisterProcessorPool.Put(opp)

	return nil
}
