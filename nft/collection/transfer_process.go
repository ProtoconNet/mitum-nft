package collection

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
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

func (op Transfer) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type TransferItemProcessor struct {
	cp           *extensioncurrency.CurrencyPool
	h            valuehash.Hash
	ns           []nft.NFT
	nftStates    map[nft.NFTID]state.State
	nftIdHStates map[nft.NFTID]state.State
	burns        []nft.NFTID
	sender       base.Address
	item         TransferItem
}

func (ipp *TransferItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return operation.NewBaseReasonError(err.Error())
	}

	if !ipp.item.Receiver().Equal(nft.BLACKHOLE_ZERO) {
		if err := checkExistsState(currency.StateKeyAccount(ipp.item.Receiver()), getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		}
		if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(ipp.item.Receiver()), getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		}
	} else {
		ipp.burns = append(ipp.burns, ipp.item.NFTs()...)
	}

	var n nft.NFT
	var nftState state.State
	nfts := ipp.item.NFTs()
	for i := range nfts {
		if err := nfts[i].IsValid(nil); err != nil {
			return operation.NewBaseReasonError(err.Error())
		}

		if st, err := existsState(StateKeyNFT(nfts[i]), "nft", getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else if _n, err := StateNFTValue(st); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else {
			n = _n
			nftState = st
		}

		if st, err := existsState(StateKeyIDFromNFTHash(n.NftHash().String()), "nft hash", getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else if id, err := StateIDFromNFTHashValue(st); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else if err := id.IsValid(nil); err != nil {
			return operation.NewBaseReasonError("dead nft hash; %q", n.NftHash())
		} else {
			ipp.nftIdHStates[n.ID()] = st
		}

		if st, err := existsState(StateKeyCollection(n.ID().Collection()), "collection", getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else if design, err := StateCollectionValue(st); err != nil {
			return operation.NewBaseReasonError(err.Error())
		} else if !design.Active() {
			return operation.NewBaseReasonError("deactivated collection; %q", n.ID())
		}

		if err := checkExistsState(currency.StateKeyAccount(n.Owner()), getState); err != nil {
			return operation.NewBaseReasonError(err.Error())
		}

		if !(n.Owner().Equal(ipp.sender) || n.Approved().Equal(ipp.sender)) {
			if st, err := existsState(StateKeyAgents(n.Owner()), "agents", getState); err != nil {
				return operation.NewBaseReasonError("unathorized sender; %q", ipp.sender)
			} else if box, err := StateAgentsValue(st); err != nil {
				return operation.NewBaseReasonError(err.Error())
			} else if !box.Exists(ipp.sender) {
				return operation.NewBaseReasonError("unathorized sender; %q", ipp.sender)
			}
		}

		ipp.ns = append(ipp.ns, n)
		ipp.nftStates[n.ID()] = nftState
	}

	return nil
}

func (ipp *TransferItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	ns := []nft.NFT{}
	for i := range ipp.ns {
		n := nft.NewNFT(ipp.ns[i].ID(), ipp.item.Receiver(), ipp.ns[i].NftHash(), ipp.ns[i].Uri(), nft.BLACKHOLE_ZERO, ipp.ns[i].Copyrighter())
		if err := n.IsValid(nil); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}
		ns = append(ns, n)
	}
	ipp.ns = ns

	for i := range ipp.ns {
		if st, err := SetStateNFTValue(ipp.nftStates[ipp.ns[i].ID()], ipp.ns[i]); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		} else {
			states = append(states, st)
		}
	}

	for i := range ipp.ns {
		if ipp.item.Receiver().Equal(nft.BLACKHOLE_ZERO) {
			if st, err := SetStateIDFromNFTHashValue(ipp.nftIdHStates[ipp.ns[i].ID()], nft.NFTID{}); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else {
				states = append(states, st)
			}
		} else {
			states = append(states, ipp.nftIdHStates[ipp.ns[i].ID()])
		}
	}

	return states, nil
}

func (ipp *TransferItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.ns = nil
	ipp.nftStates = nil
	ipp.nftIdHStates = nil
	ipp.sender = nil
	ipp.item = BaseTransferItem{}
	TransferItemProcessorPool.Put(ipp)

	return nil
}

type TransferProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Transfer
	boxes        map[extensioncurrency.ContractID]*NFTBox
	nboxStates   map[extensioncurrency.ContractID]state.State
	amountStates map[currency.CurrencyID]currency.AmountState
	ipps         []*TransferItemProcessor
	required     map[currency.CurrencyID][2]currency.Big
}

func NewTransferProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Transfer)
		if !ok {
			return nil, operation.NewBaseReasonError("not Transfer; %T", op)
		}

		opp := TransferProcessorPool.Get().(*TransferProcessor)

		opp.cp = cp
		opp.Transfer = i
		opp.boxes = map[extensioncurrency.ContractID]*NFTBox{}
		opp.nboxStates = map[extensioncurrency.ContractID]state.State{}
		opp.amountStates = nil
		opp.ipps = nil
		opp.required = nil

		return opp, nil

	}
}

func (opp *TransferProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(TransferFact)

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, err
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, err
	}

	ipps := make([]*TransferItemProcessor, len(fact.items))
	for i := range fact.items {

		c := TransferItemProcessorPool.Get().(*TransferItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.ns = []nft.NFT{}
		c.nftStates = map[nft.NFTID]state.State{}
		c.nftIdHStates = map[nft.NFTID]state.State{}
		c.sender = fact.Sender()
		c.item = fact.items[i]

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, err
		}

		ipps[i] = c
	}

	for i := range opp.ipps {
		nfts := opp.ipps[i].burns
		for j := range nfts {
			if st, err := existsState(StateKeyNFTs(nfts[j].Collection()), "collection nfts", getState); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else if box, err := StateNFTsValue(st); err != nil {
				return nil, operation.NewBaseReasonError(err.Error())
			} else {
				opp.nboxStates[nfts[j].Collection()] = st
				opp.boxes[nfts[j].Collection()] = &box
			}
		}
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else if sts, err := CheckSenderEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
		opp.amountStates = sts
	}

	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	opp.ipps = ipps

	return opp, nil
}

func (opp *TransferProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(TransferFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process approve item; %w", err)
		} else {
			states = append(states, sts...)
		}

		nfts := opp.ipps[i].burns
		for j := range nfts {
			if err := opp.boxes[nfts[j].Collection()].Remove(nfts[j]); err != nil {
				return operation.NewBaseReasonError(err.Error())
			}
		}
	}

	for k, v := range opp.nboxStates {
		if st, err := SetStateNFTsValue(v, *opp.boxes[k]); err != nil {
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

func (opp *TransferProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Transfer = Transfer{}
	opp.boxes = nil
	opp.nboxStates = nil
	opp.amountStates = nil
	opp.ipps = nil
	opp.required = nil

	TransferProcessorPool.Put(opp)

	return nil
}

func (opp *TransferProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(TransferFact)

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
			return nil, operation.NewBaseReasonError("unknown currency id found, %q", it.Currency())
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
