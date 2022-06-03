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
	idxState        state.State
	collectionState state.State
	nftState        state.State
	nboxState       state.State
	nftIdHStates    state.State
	idx             uint
	design          nft.Design
	n               nft.NFT
	box             NFTBox
	amountState     currency.AmountState
	fee             currency.Big
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
		opp.idxState = nil
		opp.collectionState = nil
		opp.nftState = nil
		opp.nboxState = nil
		opp.nftIdHStates = nil
		opp.idx = 0
		opp.design = nft.Design{}
		opp.n = nft.NFT{}
		opp.box = NFTBox{}
		opp.amountState = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *MintProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
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

	if st, found, _ := getState(StateKeyIDFromNFTHash(fact.Form().Hash().String())); found {
		if id, err := StateIDFromNFTHashValue(st); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		} else if err := id.IsValid(nil); err == nil {
			return nil, operation.NewBaseReasonError("nft hash %v is alive in some collection", fact.Form().Hash())
		}
		opp.nftIdHStates = st
	} else {
		opp.nftIdHStates = st
	}

	if st, err := existsState(StateKeyCollection(fact.Collection()), "design", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.collectionState = st
	}

	if design, err := StateCollectionValue(opp.collectionState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if !design.Active() {
		return nil, operation.NewBaseReasonError("deactivated collection; %q", design.Symbol())
	} else {
		opp.design = design
	}

	if !fact.Sender().Equal(opp.design.Creator()) {
		return nil, operation.NewBaseReasonError(
			"sender must be collection creator; creator: %q", opp.design.Creator().String())
	}

	if st, err := existsState(StateKeyCollectionLastIDX(fact.Collection()), "collection idx", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.idxState = st
	}

	if idx, err := StateCollectionLastIDXValue(opp.idxState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.idx = idx + 1
	}

	switch st, found, err := getState(StateKeyNFTs(fact.Collection())); {
	case err != nil:
		return nil, operation.NewBaseReasonError(err.Error())
	case !found:
		opp.box = NewNFTBox(nil)
		opp.nboxState = st
	default:
		box, err := StateNFTsValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}
		opp.box = box
		opp.nboxState = st
	}

	id := nft.NewNFTID(fact.Collection(), opp.idx)
	if err := id.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if st, err := notExistsState(StateKeyNFT(id), "nft", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.nftState = st
	}

	if fact.Form().Copyrighter().String() != "" {
		if err := checkExistsState(currency.StateKeyAccount(fact.Form().Copyrighter()), getState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		} else if err = checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Form().Copyrighter()), getState); err != nil {
			return nil, operation.NewBaseReasonError("contract account cannot be copyrighter; %q", fact.Sender())
		}
	}

	n := nft.NewNFT(id, fact.Sender(), fact.Form().Hash(), fact.Form().Uri(), nft.BLACKHOLE_ZERO, fact.Form().Copyrighter())
	if err := n.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}
	opp.n = n

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	if st, err := existsState(
		currency.StateKeyBalance(fact.Sender(), fact.Currency()), "balance of sender", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
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

func (opp *MintProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(MintFact)

	var states []state.State

	if st, err := SetStateCollectionLastIDXValue(opp.idxState, opp.idx); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	if st, err := SetStateNFTValue(opp.nftState, opp.n); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	if st, err := SetStateIDFromNFTHashValue(opp.nftIdHStates, opp.n.ID()); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	if err := opp.box.Append(opp.n.ID()); err != nil {
		return operation.NewBaseReasonError(err.Error())
	}
	if st, err := SetStateNFTsValue(opp.nboxState, opp.box); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	opp.amountState = opp.amountState.Sub(opp.fee).AddFee(opp.fee)
	states = append(states, opp.amountState)

	return setState(fact.Hash(), states...)
}

func (opp *MintProcessor) Close() error {
	opp.cp = nil
	opp.Mint = Mint{}
	opp.idxState = nil
	opp.collectionState = nil
	opp.nftState = nil
	opp.nboxState = nil
	opp.nftIdHStates = nil
	opp.idx = 0
	opp.design = nft.Design{}
	opp.n = nft.NFT{}
	opp.box = NFTBox{}
	opp.amountState = currency.AmountState{}
	opp.fee = currency.ZeroBig

	MintProcessorPool.Put(opp)

	return nil
}
