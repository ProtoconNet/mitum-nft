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
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	box    *NFTBox
	nft    nft.NFT
	nst    state.State
	sender base.Address
	item   BurnItem
}

func (ipp *BurnItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	nid := ipp.item.NFT()

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

		n := nft.NewNFT(nv.ID(), false, nv.Owner(), nv.NftHash(), nv.Uri(), nv.Owner(), nv.Creators(), nv.Copyrighters())
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

func (ipp *BurnItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	if err := ipp.box.Remove(ipp.nft.ID()); err != nil {
		return nil, err
	}

	if st, err := SetStateNFTValue(ipp.nst, ipp.nft); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *BurnItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.nft = nft.NFT{}
	ipp.nst = nil
	ipp.box = nil
	ipp.sender = nil
	ipp.item = BurnItem{}
	BurnItemProcessorPool.Put(ipp)

	return nil
}

type BurnProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Burn
	boxes        map[extensioncurrency.ContractID]*NFTBox
	boxStates    map[extensioncurrency.ContractID]state.State
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
		opp.boxes = nil
		opp.boxStates = nil
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
	fact, ok := opp.Fact().(BurnFact)
	if !ok {
		return nil, operation.NewBaseReasonError("not BurnFact; %T", opp.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot burn nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	opp.boxes = map[extensioncurrency.ContractID]*NFTBox{}
	opp.boxStates = map[extensioncurrency.ContractID]state.State{}
	for i := range fact.items {
		collection := fact.items[i].NFT().Collection()

		if _, found := opp.boxes[collection]; !found {
			var box NFTBox
			if st, err := existsState(StateKeyCollection(collection), "design", getState); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if design, err := StateCollectionValue(st); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if !design.Active() {
				return nil, operation.NewBaseReasonError("deactivated collection; %q", collection)
			}

			if st, err := existsState(StateKeyNFTs(collection), "nfts", getState); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if b, err := StateNFTsValue(st); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else {
				box = b
				opp.boxStates[collection] = st
			}
			opp.boxes[collection] = &box
		}

	}

	ipps := make([]*BurnItemProcessor, len(fact.items))
	for i := range fact.items {

		c := BurnItemProcessorPool.Get().(*BurnItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.box = opp.boxes[fact.items[i].NFT().Collection()]
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

func (opp *BurnProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact, ok := opp.Fact().(BurnFact)
	if !ok {
		return operation.NewBaseReasonError("not BurnFact; %T", opp.Fact())
	}

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process burn item; %w", err)
		} else {
			states = append(states, sts...)
		}
	}

	for c, box := range opp.boxes {
		if st, err := SetStateNFTsValue(opp.boxStates[c], *box); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else {
			states = append(states, st)
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
	opp.boxes = nil
	opp.boxStates = nil
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil

	BurnProcessorPool.Put(opp)

	return nil
}

func (opp *BurnProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := opp.Fact().(BurnFact)
	if !ok {
		return nil, errors.Errorf("not BurnFact; %T", opp.Fact())
	}

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
			return nil, errors.Errorf("unknown currency id found; %q", it.Currency())
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
