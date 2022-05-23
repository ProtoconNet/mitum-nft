package collection

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	DelegateFactType   = hint.Type("mitum-nft-delegate-operation-fact")
	DelegateFactHint   = hint.NewHint(DelegateFactType, "v0.0.1")
	DelegateFactHinter = DelegateFact{BaseHinter: hint.NewBaseHinter(DelegateFactHint)}
	DelegateType       = hint.Type("mitum-nft-delegate-operation")
	DelegateHint       = hint.NewHint(DelegateType, "v0.0.1")
	DelegateHinter     = Delegate{BaseOperation: operationHinter(DelegateHint)}
)

var MaxAgents = 10

type DelegateFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []DelegateItem
}

func NewDelegateFact(token []byte, sender base.Address, items []DelegateItem) DelegateFact {
	fact := DelegateFact{
		BaseHinter: hint.NewBaseHinter(DelegateFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact DelegateFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact DelegateFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DelegateFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact DelegateFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for DelegateFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender); err != nil {
		return err
	}

	foundAgent := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		agent := fact.items[i].Agent()
		if err := agent.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundAgent[agent.String()]; found {
			return isvalid.InvalidError.Errorf("duplicated agent found; %s", agent)
		}
		foundAgent[agent.String()] = true
	}

	return nil
}

func (fact DelegateFact) Token() []byte {
	return fact.token
}

func (fact DelegateFact) Sender() base.Address {
	return fact.sender
}

func (fact DelegateFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	for i := range fact.items {
		as[i] = fact.items[i].Agent()
	}

	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact DelegateFact) Currencies() []currency.CurrencyID {
	cs := make([]currency.CurrencyID, len(fact.items))

	for i := range fact.items {
		cs[i] = fact.items[i].Currency()
	}

	return cs
}

func (fact DelegateFact) Rebuild() DelegateFact {
	items := make([]DelegateItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Delegate struct {
	currency.BaseOperation
}

func NewDelegate(fact DelegateFact, fs []base.FactSign, memo string) (Delegate, error) {
	bo, err := currency.NewBaseOperationFromFact(DelegateHint, fact, fs, memo)
	if err != nil {
		return Delegate{}, err
	}
	return Delegate{BaseOperation: bo}, nil
}
