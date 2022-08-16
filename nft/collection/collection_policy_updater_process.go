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

var CollectionPolicyUpdaterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionPolicyUpdaterProcessor)
	},
}

func (CollectionPolicyUpdater) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CollectionPolicyUpdaterProcessor struct {
	cp *extensioncurrency.CurrencyPool
	CollectionPolicyUpdater
	designState state.State
	design      nft.Design
	amountState currency.AmountState
	fee         currency.Big
}

func NewCollectionPolicyUpdaterProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(CollectionPolicyUpdater)
		if !ok {
			return nil, errors.Errorf("not CollectionPolicyUpdater; %T", op)
		}

		opp := CollectionPolicyUpdaterProcessorPool.Get().(*CollectionPolicyUpdaterProcessor)

		opp.cp = cp
		opp.CollectionPolicyUpdater = i
		opp.designState = nil
		opp.design = nft.Design{}
		opp.amountState = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *CollectionPolicyUpdaterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact, ok := opp.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return nil, operation.NewBaseReasonError("not CollectionPolicyUpdaterFact; %T", opp.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot update collection policy; %q", fact.Sender())
	}

	if st, err := existsState(StateKeyCollection(fact.Collection()), "design", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if design, err := StateCollectionValue(st); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if !design.Active() {
		return nil, operation.NewBaseReasonError("deactivated collection; %q", fact.Collection())
	} else if !design.Creator().Equal(fact.Sender()) {
		return nil, operation.NewBaseReasonError("not creator of collection design; %q", fact.Collection())
	} else if cst, err := existsState(extensioncurrency.StateKeyContractAccount(design.Parent()), "contract account", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if ca, err := extensioncurrency.StateContractAccountValue(cst); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if !ca.IsActive() {
		return nil, operation.NewBaseReasonError("deactivated contract account; %q", design.Parent())
	} else {
		opp.designState = st
		opp.design = nft.NewDesign(design.Parent(), design.Creator(), design.Symbol(), design.Active(), fact.Policy())
	}

	if err := opp.design.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

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

func (opp *CollectionPolicyUpdaterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact, ok := opp.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return operation.NewBaseReasonError("not CollectionPolicyUpdaterFact; %T", opp.Fact())
	}

	var states []state.State

	if st, err := SetStateCollectionValue(opp.designState, opp.design); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	opp.amountState = opp.amountState.Sub(opp.fee).AddFee(opp.fee)
	states = append(states, opp.amountState)

	return setState(fact.Hash(), states...)
}

func (opp *CollectionPolicyUpdaterProcessor) Close() error {
	opp.cp = nil
	opp.designState = nil
	opp.design = nft.Design{}
	opp.amountState = currency.AmountState{}
	opp.fee = currency.ZeroBig

	CollectionPolicyUpdaterProcessorPool.Put(opp)

	return nil
}
