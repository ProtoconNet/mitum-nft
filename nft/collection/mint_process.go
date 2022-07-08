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

var MintItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintItemProcessor)
	},
}

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

type MintItemProcessor struct {
	cp       *extensioncurrency.CurrencyPool
	h        valuehash.Hash
	idxState state.State
	idx      uint64
	boxState state.State
	box      NFTBox
	nft      nft.NFT
	nst      state.State
	sender   base.Address
	item     MintItem
}

func (ipp *MintItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if st, err := existsState(StateKeyCollection(ipp.item.Collection()), "design", getState); err != nil {
		return err
	} else if design, err := StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("deactivated collection; %q", design.Symbol())
	} else if !ipp.sender.Equal(design.Creator()) {
		return errors.Errorf("sender must be collection creator; creator: %q", design.Creator().String())
	}

	if st, err := existsState(StateKeyCollectionLastIDX(ipp.item.Collection()), "collection idx", getState); err != nil {
		return err
	} else if idx, err := StateCollectionLastIDXValue(st); err != nil {
		return err
	} else {
		ipp.idxState = st
		ipp.idx = idx + 1
	}

	switch st, found, err := getState(StateKeyNFTs(ipp.item.Collection())); {
	case err != nil:
		return err
	case !found:
		ipp.box = NewNFTBox(nil)
		ipp.boxState = st
	default:
		box, err := StateNFTsValue(st)
		if err != nil {
			return err
		}
		ipp.box = box
		ipp.boxState = st
	}

	id := nft.NewNFTID(ipp.item.Collection(), ipp.idx)
	if err := id.IsValid(nil); err != nil {
		return err
	}

	if st, err := notExistsState(StateKeyNFT(id), "nft", getState); err != nil {
		return err
	} else {
		ipp.nst = st
	}

	form := ipp.item.Form()
	if form.Creators().Total() != 0 {
		creators := form.Creators().Signers()
		for i := range creators {
			creator := creators[i].Account()
			if err := checkExistsState(currency.StateKeyAccount(creator), getState); err != nil {
				return err
			} else if err = checkNotExistsState(extensioncurrency.StateKeyContractAccount(creator), getState); err != nil {
				return errors.Errorf("contract account cannot be a creator; %q", creator)
			}
			if creators[i].Signed() {
				return errors.Errorf("Cannot sign at the same time as minting; %q", creator)
			}
		}
	}

	if form.Copyrighters().Total() != 0 {
		copyrighters := form.Copyrighters().Signers()
		for i := range copyrighters {
			copyrighter := copyrighters[i].Account()
			if err := checkExistsState(currency.StateKeyAccount(copyrighter), getState); err != nil {
				return err
			} else if err = checkNotExistsState(extensioncurrency.StateKeyContractAccount(copyrighter), getState); err != nil {
				return errors.Errorf("contract account cannot be a copyrighter; %q", copyrighter)
			}
			if copyrighters[i].Signed() {
				return errors.Errorf("Cannot sign at the same time as minting; %q", copyrighter)
			}
		}
	}

	n := nft.NewNFT(id, true, ipp.sender, form.NftHash(), form.Uri(), ipp.sender, form.Creators(), form.Copyrighters())
	if err := n.IsValid(nil); err != nil {
		return operation.NewBaseReasonError(err.Error())
	}
	ipp.nft = n

	return nil
}

func (ipp *MintItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	if st, err := SetStateCollectionLastIDXValue(ipp.idxState, ipp.idx); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	if err := ipp.box.Append(ipp.nft.ID()); err != nil {
		return nil, err
	}

	if st, err := SetStateNFTValue(ipp.nst, ipp.nft); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	if st, err := SetStateNFTsValue(ipp.boxState, ipp.box); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *MintItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.idxState = nil
	ipp.idx = 0
	ipp.boxState = nil
	ipp.box = NFTBox{}
	ipp.nft = nft.NFT{}
	ipp.nst = nil
	ipp.sender = nil
	ipp.item = MintItem{}
	MintItemProcessorPool.Put(ipp)

	return nil
}

type MintProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Mint
	ipps         []*MintItemProcessor
	amountStates map[currency.CurrencyID]currency.AmountState
	required     map[currency.CurrencyID][2]currency.Big
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
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *MintProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(MintFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot mint nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	ipps := make([]*MintItemProcessor, len(fact.items))
	for i := range fact.items {

		c := MintItemProcessorPool.Get().(*MintItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.idxState = nil
		c.idx = 0
		c.boxState = nil
		c.box = NFTBox{}
		c.nft = nft.NFT{}
		c.nst = nil
		c.sender = fact.Sender()
		c.item = fact.items[i]

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}

		ipps[i] = c
	}

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

	opp.ipps = ipps

	return opp, nil
}

func (opp *MintProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(MintFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process mint item; %w", err)
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

func (opp *MintProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Mint = Mint{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil

	MintProcessorPool.Put(opp)

	return nil
}

func (opp *MintProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(MintFact)

	items := make([]MintItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateMintItemsFee(opp.cp, items)
}

func CalculateMintItemsFee(cp *extensioncurrency.CurrencyPool, items []MintItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
