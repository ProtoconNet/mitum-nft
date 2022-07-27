package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var TransferItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferItemProcessor)
	},
}

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

type TransferItemProcessor struct {
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	nft    nft.NFT
	nst    state.State
	sender base.Address
	item   TransferItem
}

func (ipp *TransferItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	// check receiver
	receiver := ipp.item.Receiver()
	if err := checkExistsState(currency.StateKeyAccount(receiver), getState); err != nil {
		return err
	}
	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(receiver), getState); err != nil {
		return errors.Errorf("contract account cannot receive nfts; %q", receiver)
	}

	nid := ipp.item.NFT()
	if st, err := existsState(StateKeyCollection(nid.Collection()), "design", getState); err != nil {
		return errors.Errorf("%v; %q", err.Error(), nid.Collection())
	} else if design, err := StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("deactivated collection; %q", design.Symbol())
	}

	var (
		approved base.Address
		owner    base.Address
	)

	// check nft
	if st, err := existsState(StateKeyNFT(nid), "nft", getState); err != nil {
		return err
	} else if nv, err := StateNFTValue(st); err != nil {
		return err
	} else if !nv.Active() {
		return errors.Errorf("burned nft; %q", nid)
	} else {
		approved = nv.Approved()
		owner = nv.Owner()

		n := nft.NewNFT(nid, nv.Active(), receiver, nv.NftHash(), nv.Uri(), receiver, nv.Creators(), nv.Copyrighters())
		if err := n.IsValid(nil); err != nil {
			return err
		}

		ipp.nft = n
		ipp.nst = st
	}

	// check authorization
	if !(owner.Equal(ipp.sender) || approved.Equal(ipp.sender)) {
		// check agent
		if st, err := existsState(StateKeyAgents(owner, ipp.nft.ID().Collection()), "agents", getState); err != nil {
			return errors.Errorf("unauthorized sender; %q", ipp.sender)
		} else if box, err := StateAgentsValue(st); err != nil {
			return err
		} else if !box.Exists(ipp.sender) {
			return errors.Errorf("unauthorized sender; %q", ipp.sender)
		}
	}

	return nil
}

func (ipp *TransferItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	if st, err := SetStateNFTValue(ipp.nst, ipp.nft); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *TransferItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.nft = nft.NFT{}
	ipp.nst = nil
	ipp.sender = nil
	ipp.item = TransferItem{}
	TransferItemProcessorPool.Put(ipp)

	return nil
}

type TransferProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Transfer
	ipps         []*TransferItemProcessor
	amountStates map[currency.CurrencyID]currency.AmountState
	required     map[currency.CurrencyID][2]currency.Big
}

func NewTransferProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Transfer)
		if !ok {
			return nil, errors.Errorf("not Transfer; %T", op)
		}

		opp := TransferProcessorPool.Get().(*TransferProcessor)

		opp.cp = cp
		opp.Transfer = i
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *TransferProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact, ok := opp.Fact().(TransferFact)
	if !ok {
		return nil, operation.NewBaseReasonError("not TransferFact; %T", opp.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot transfer nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	ipps := make([]*TransferItemProcessor, len(fact.items))
	for i := range fact.items {

		c := TransferItemProcessorPool.Get().(*TransferItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.nft = nft.NFT{}
		c.nst = nil
		c.sender = fact.Sender()
		c.item = fact.items[i]

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}

		ipps[i] = c
	}

	opp.ipps = ipps

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else if sts, err := CheckSenderEnoughBalance(fact.Sender(), required, getState); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else {
		opp.required = required
		opp.amountStates = sts
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	return opp, nil
}

func (opp *TransferProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact, ok := opp.Fact().(TransferFact)
	if !ok {
		return operation.NewBaseReasonError("not TransferFact; %T", opp.Fact())
	}

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process transfer item; %w", err)
		} else {
			states = append(states, sts...)
		}
	}

	for k := range opp.required {
		rq := opp.required[k]
		states = append(states, opp.amountStates[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), states...)
}

func (opp *TransferProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Transfer = Transfer{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil

	TransferProcessorPool.Put(opp)

	return nil
}

func (opp *TransferProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := opp.Fact().(TransferFact)
	if !ok {
		return nil, errors.Errorf("not TransferFact; %T", opp.Fact())
	}

	items := make([]TransferItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateTransferItemsFee(opp.cp, items)
}

func CalculateTransferItemsFee(cp *extensioncurrency.CurrencyPool, items []TransferItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
