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

var ApproveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var ApproveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type ApproveItemProcessor struct {
	cp     *extensioncurrency.CurrencyPool
	h      valuehash.Hash
	nft    nft.NFT
	nst    state.State
	sender base.Address
	item   ApproveItem
}

func (ipp *ApproveItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if ipp.item.Approved().String() != "" {
		if err := checkExistsState(currency.StateKeyAccount(ipp.item.Approved()), getState); err != nil {
			return err
		}
		if ipp.item.Approved().Equal(ipp.sender) {
			return errors.Errorf("sender cannot be approved account itself; %q", ipp.item.Approved())
		}
	}

	nid := ipp.item.NFT()

	if st, err := existsState(StateKeyCollection(nid.Collection()), "design", getState); err != nil {
		return err
	} else if design, err := StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("dead collection; %q", nid.Collection())
	}

	if st, err := existsState(StateKeyNFT(nid), "nft", getState); err != nil {
		return err
	} else if nv, err := StateNFTValue(st); err != nil {
		return err
	} else {
		ipp.nft = nv
		ipp.nst = st
	}

	if ipp.nft.Owner().String() == "" {
		return errors.Errorf("dead nft; %q", nid)
	}

	if !ipp.nft.Owner().Equal(ipp.sender) {
		if err := checkExistsState(currency.StateKeyAccount(ipp.nft.Owner()), getState); err != nil {
			return err
		} else if st, err := existsState(StateKeyAgents(ipp.nft.Owner(), ipp.nft.ID().Collection()), "agents", getState); err != nil {
			return errors.Errorf("unathorized sender; %q", ipp.sender)
		} else if box, err := StateAgentsValue(st); err != nil {
			return err
		} else if !box.Exists(ipp.sender) {
			return errors.Errorf("unathorized sender; %q", ipp.sender)
		}
	}

	n := nft.NewNFT(
		nid, ipp.nft.Owner(), ipp.nft.NftHash(),
		ipp.nft.Uri(), ipp.item.Approved(), ipp.nft.Creators(), ipp.nft.Copyrighters(),
	)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	ipp.nft = n

	return nil
}

func (ipp *ApproveItemProcessor) Process(
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

func (ipp *ApproveItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.nft = nft.NFT{}
	ipp.nst = nil
	ipp.sender = nil
	ipp.item = ApproveItem{}
	ApproveItemProcessorPool.Put(ipp)

	return nil
}

type ApproveProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Approve
	amountStates map[currency.CurrencyID]currency.AmountState
	ipps         []*ApproveItemProcessor
	required     map[currency.CurrencyID][2]currency.Big
}

func NewApproveProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Approve)
		if !ok {
			return nil, operation.NewBaseReasonError("not Approve; %T", op)
		}

		opp := ApproveProcessorPool.Get().(*ApproveProcessor)

		opp.cp = cp
		opp.Approve = i
		opp.amountStates = nil
		opp.ipps = nil
		opp.required = nil

		return opp, nil

	}
}

func (opp *ApproveProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(ApproveFact)

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	ipps := make([]*ApproveItemProcessor, len(fact.items))
	for i := range fact.items {

		c := ApproveItemProcessorPool.Get().(*ApproveItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.sender = fact.Sender()
		c.item = fact.items[i]
		c.nft = nft.NFT{}
		c.nst = nil

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}

		ipps[i] = c
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else if sts, err := CheckSenderEnoughBalance(fact.Sender(), required, getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
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

func (opp *ApproveProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(ApproveFact)

	var states []state.State

	for i := range opp.ipps {
		if s, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process approve item; %w", err)
		} else {
			states = append(states, s...)
		}
	}

	for k := range opp.required {
		rq := opp.required[k]
		states = append(states, opp.amountStates[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), states...)
}

func (opp *ApproveProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Approve = Approve{}
	opp.amountStates = nil
	opp.required = nil
	opp.ipps = nil

	ApproveProcessorPool.Put(opp)

	return nil
}

func (opp *ApproveProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(ApproveFact)

	items := make([]ApproveItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateApproveItemsFee(opp.cp, items)
}

func CalculateApproveItemsFee(cp *extensioncurrency.CurrencyPool, items []ApproveItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
