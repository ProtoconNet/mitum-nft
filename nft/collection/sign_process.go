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

var SignItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignItemProcessor)
	},
}

var SignProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignProcessor)
	},
}

func (Sign) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type SignItemProcessor struct {
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	nft    nft.NFT
	nst    state.State
	sender base.Address
	item   SignItem
}

func (ipp *SignItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	nid := ipp.item.NFT()

	// check collection
	if st, err := existsState(StateKeyCollection(nid.Collection()), "design", getState); err != nil {
		return err
	} else if design, err := StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("dead collection; %q", nid.Collection())
	}

	var signers []nft.Signer
	var n nft.NFT

	// check nft
	if st, err := existsState(StateKeyNFT(nid), "nft", getState); err != nil {
		return err
	} else if nv, err := StateNFTValue(st); err != nil {
		return err
	} else {
		switch ipp.item.Qualification() {
		case CreatorQualification:
			signers = nv.Creators()
		case CopyrighterQualification:
			signers = nv.Copyrighters()
		default:
			return errors.Errorf("wrong qualification; %q", ipp.item.Qualification())
		}
		n = nv
		ipp.nst = st
	}

	// check owner
	if n.Owner().String() == "" {
		return errors.Errorf("dead nft; %q", nid)
	}

	var idx = -1
	for i := range signers {
		if signers[i].Account().Equal(ipp.sender) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return errors.Errorf("not signer of nft; %q, %q", ipp.sender, n.ID())
	}

	if signers[idx].Signed() {
		return errors.Errorf("this signer has already signed nft; %q", signers[idx].Account())
	}

	signer := nft.NewSigner(ipp.sender, true)
	if err := signer.IsValid(nil); err != nil {
		return err
	}
	signers[idx] = signer

	if ipp.item.Qualification() == CreatorQualification {
		n = nft.NewNFT(n.ID(), n.Owner(), n.NftHash(), n.Uri(), n.Approved(), signers, n.Copyrighters())
	} else {
		n = nft.NewNFT(n.ID(), n.Owner(), n.NftHash(), n.Uri(), n.Approved(), n.Creators(), signers)
	}

	if err := n.IsValid(nil); err != nil {
		return err
	}
	ipp.nft = n

	return nil
}

func (ipp *SignItemProcessor) Process(
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

func (ipp *SignItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.nft = nft.NFT{}
	ipp.nst = nil
	ipp.sender = nil
	ipp.item = SignItem{}
	SignItemProcessorPool.Put(ipp)

	return nil
}

type SignProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Sign
	ipps         []*SignItemProcessor
	amountStates map[currency.CurrencyID]currency.AmountState
	required     map[currency.CurrencyID][2]currency.Big
}

func NewSignProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Sign)
		if !ok {
			return nil, errors.Errorf("not Sign; %T", op)
		}

		opp := SignProcessorPool.Get().(*SignProcessor)

		opp.cp = cp
		opp.Sign = i
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil

		return opp, nil
	}
}

func (opp *SignProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(SignFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot sign nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	ipps := make([]*SignItemProcessor, len(fact.items))
	for i := range fact.items {

		c := SignItemProcessorPool.Get().(*SignItemProcessor)
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

func (opp *SignProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(SignFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process sign item; %w", err)
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

func (opp *SignProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Sign = Sign{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil

	SignProcessorPool.Put(opp)

	return nil
}

func (opp *SignProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(SignFact)

	items := make([]SignItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateSignItemsFee(opp.cp, items)
}

func CalculateSignItemsFee(cp *extensioncurrency.CurrencyPool, items []SignItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
