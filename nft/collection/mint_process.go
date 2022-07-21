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
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	idx    uint64
	box    *NFTBox
	nft    nft.NFT
	nst    state.State
	sender base.Address
	item   MintItem
}

func (ipp *MintItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if err := ipp.item.IsValid(nil); err != nil {
		return err
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

	if err := ipp.box.Append(ipp.nft.ID()); err != nil {
		return nil, err
	}

	if st, err := SetStateNFTValue(ipp.nst, ipp.nft); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *MintItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.idx = 0
	ipp.box = nil
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
	idxes        map[extensioncurrency.ContractID]uint64
	idxStates    map[extensioncurrency.ContractID]state.State
	boxes        map[extensioncurrency.ContractID]*NFTBox
	boxStates    map[extensioncurrency.ContractID]state.State
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
		opp.idxes = nil
		opp.idxStates = nil
		opp.boxes = nil
		opp.boxStates = nil
		opp.amountStates = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *MintProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact, ok := opp.Fact().(MintFact)
	if !ok {
		return nil, operation.NewBaseReasonError("not MintFact; %T", opp.Fact())
	}

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

	opp.idxes = map[extensioncurrency.ContractID]uint64{}
	opp.idxStates = map[extensioncurrency.ContractID]state.State{}
	opp.boxes = map[extensioncurrency.ContractID]*NFTBox{}
	opp.boxStates = map[extensioncurrency.ContractID]state.State{}
	for i := range fact.items {
		collection := fact.items[i].Collection()

		if _, found := opp.idxes[collection]; !found {
			if st, err := existsState(StateKeyCollection(collection), "design", getState); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if design, err := StateCollectionValue(st); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if !design.Active() {
				return nil, operation.NewBaseReasonError("deactivated collection; %q", design.Symbol())
			} else if policy, ok := design.Policy().(CollectionPolicy); !ok {
				return nil, operation.NewBaseReasonError("policy of design is not collection-policy; %q", design.Symbol())
			} else if whites := policy.Whites(); len(whites) == 0 {
				return nil, operation.NewBaseReasonError("empty whitelist! nobody can mint to this collection; %q", collection)
			} else {
				for i := range whites {
					if whites[i].Equal(fact.Sender()) {
						break
					}
					if i == len(whites)-1 {
						return nil, operation.NewBaseReasonError("sender is not whitelisted; %q", fact.Sender())
					}
				}
			}

			if st, err := existsState(StateKeyCollectionLastIDX(collection), "collection idx", getState); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if idx, err := StateCollectionLastIDXValue(st); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else {
				opp.idxes[collection] = idx
				opp.idxStates[collection] = st
			}
		}

		if _, found := opp.boxes[collection]; !found {
			var box NFTBox
			switch st, found, err := getState(StateKeyNFTs(collection)); {
			case err != nil:
				return nil, operation.NewBaseReasonError(err.Error())
			case !found:
				box = NewNFTBox(nil)
				opp.boxStates[collection] = st
			default:
				b, err := StateNFTsValue(st)
				if err != nil {
					return nil, operation.NewBaseReasonError(err.Error())
				}
				box = b
				opp.boxStates[collection] = st
			}
			opp.boxes[collection] = &box
		}
	}

	ipps := make([]*MintItemProcessor, len(fact.items))
	for i := range fact.items {
		collection := fact.items[i].Collection()
		idx := opp.idxes[collection] + 1
		opp.idxes[collection] = idx

		c := MintItemProcessorPool.Get().(*MintItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.idx = idx
		c.box = opp.boxes[collection]
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
	fact, ok := opp.Fact().(MintFact)
	if !ok {
		return operation.NewBaseReasonError("not MintFact; %T", opp.Fact())
	}

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process mint item; %w", err)
		} else {
			states = append(states, sts...)
		}
	}

	for c, idx := range opp.idxes {
		if st, err := SetStateCollectionLastIDXValue(opp.idxStates[c], idx); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else {
			states = append(states, st)
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

func (opp *MintProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Mint = Mint{}
	opp.ipps = nil
	opp.idxes = nil
	opp.idxStates = nil
	opp.boxes = nil
	opp.boxStates = nil
	opp.amountStates = nil
	opp.required = nil

	MintProcessorPool.Put(opp)

	return nil
}

func (opp *MintProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := opp.Fact().(MintFact)
	if !ok {
		return nil, errors.Errorf("not MintFact; %T", opp.Fact())
	}

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
