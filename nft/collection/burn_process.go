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

var BurnItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(BurnItemProcessor)
	},
}

var BurnProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(BurnProcessor)
	},
}

func (Burn) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type BurnItemProcessor struct {
	cp       *extensioncurrency.CurrencyPool
	h        valuehash.Hash
	box      NFTBox
	boxState state.State
	ns       map[nft.NFT]state.State
	sender   base.Address
	item     BurnItem
}

func (ipp *BurnItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	// check collection
	if st, err := existsState(StateKeyCollection(ipp.item.Collection()), "design", getState); err != nil {
		return err
	} else if design, err := StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("something wrong~ this collection has already been deactivated; %q", ipp.item.Collection())
	}

	if st, err := existsState(StateKeyNFTs(ipp.item.Collection()), "nfts", getState); err != nil {
		return err
	} else if box, err := StateNFTsValue(st); err != nil {
		return err
	} else {
		ipp.box = box
		ipp.boxState = st
	}

	nfts := ipp.item.NFTs()
	for i := range nfts {
		var n nft.NFT
		var (
			approved base.Address
			owner    base.Address
		)

		// check nft
		if st, err := existsState(StateKeyNFT(nfts[i]), "nft", getState); err != nil {
			return err
		} else if _n, err := StateNFTValue(st); err != nil {
			return err
		} else if _n.ID().Collection() != ipp.item.Collection() {
			return errors.Errorf("another collection found; %q", _n.ID())
		} else {
			approved = _n.Approved()
			owner = _n.Owner()

			n := nft.NewNFT(_n.ID(), currency.Address{}, _n.NftHash(), _n.Uri(), currency.Address{}, _n.Copyrighter())
			if err := n.IsValid(nil); err != nil {
				return err
			}
			ipp.ns[n] = st
		}

		// check owner
		if owner.String() == "" {
			return errors.Errorf("no owner in nft; %q", n.ID())
		}
		if err := checkExistsState(currency.StateKeyAccount(owner), getState); err != nil {
			return err
		}
		if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(owner), getState); err != nil {
			return err
		}

		// check authorization
		if !(owner.Equal(ipp.sender) || approved.Equal(ipp.sender)) {
			// check agent
			if st, err := existsState(StateKeyAgents(owner), "agents", getState); err != nil {
				return errors.Errorf("unauthorized sender; %q", ipp.sender)
			} else if box, err := StateAgentsValue(st); err != nil {
				return err
			} else if !box.Exists(ipp.sender) {
				return errors.Errorf("unauthorized sender; %q", ipp.sender)
			}
		}
	}

	return nil
}

func (ipp *BurnItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	// set nfts
	for k, v := range ipp.ns {
		if st, err := SetStateNFTValue(v, k); err != nil {
			return nil, err
		} else {
			states = append(states, st)
		}
	}

	// remove nfts
	nfts := ipp.item.NFTs()
	for i := range ipp.item.NFTs() {
		if err := ipp.box.Remove(nfts[i]); err != nil {
			return nil, err
		}
	}

	// set nft box
	if st, err := SetStateNFTsValue(ipp.boxState, ipp.box); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *BurnItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.ns = nil
	ipp.box = NFTBox{}
	ipp.boxState = nil
	ipp.sender = nil
	ipp.item = BaseBurnItem{}
	BurnItemProcessorPool.Put(ipp)

	return nil
}

type BurnProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Burn
	ipps         []*BurnItemProcessor
	amountStates map[currency.CurrencyID]currency.AmountState
	required     map[currency.CurrencyID][2]currency.Big
}

func NewBurnProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Burn)
		if !ok {
			return nil, errors.Errorf("not Burn; %T", op)
		}

		opp := BurnProcessorPool.Get().(*BurnProcessor)

		opp.cp = cp
		opp.Burn = i
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *BurnProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(BurnFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot Burn nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	ipps := make([]*BurnItemProcessor, len(fact.items))
	for i := range fact.items {

		c := BurnItemProcessorPool.Get().(*BurnItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.box = NFTBox{}
		c.boxState = nil
		c.ns = map[nft.NFT]state.State{}
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

func (opp *BurnProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(BurnFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process Burn item; %w", err)
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

func (opp *BurnProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Burn = Burn{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil

	BurnProcessorPool.Put(opp)

	return nil
}

func (opp *BurnProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(BurnFact)

	items := make([]BurnItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateBurnItemsFee(opp.cp, items)
}

func CalculateBurnItemsFee(cp *extensioncurrency.CurrencyPool, items []BurnItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
