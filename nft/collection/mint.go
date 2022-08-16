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

var MaxMintItems = 10

var (
	MintFactType   = hint.Type("mitum-nft-mint-operation-fact")
	MintFactHint   = hint.NewHint(MintFactType, "v0.0.1")
	MintFactHinter = MintFact{BaseHinter: hint.NewBaseHinter(MintFactHint)}
	MintType       = hint.Type("mitum-nft-mint-operation")
	MintHint       = hint.NewHint(MintType, "v0.0.1")
	MintHinter     = Mint{BaseOperation: operationHinter(MintHint)}
)

type MintFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []MintItem
}

func NewMintFact(token []byte, sender base.Address, items []MintItem) MintFact {
	fact := MintFact{
		BaseHinter: hint.NewBaseHinter(MintFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact MintFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact MintFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
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

func (fact MintFact) IsValid(b []byte) error {
	if err := isvalid.Check(
		nil, false,
		fact.BaseHinter,
		fact.h,
		fact.sender); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for MintFact")
	}

	if l := len(fact.items); l < 1 {
		return isvalid.InvalidError.Errorf("empty items for MintFact")
	} else if l > int(MaxMintItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", l, MaxMintItems)
	}

	for i := range fact.items {
		if err := fact.items[i].IsValid(nil); err != nil {
			return err
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact MintFact) Token() []byte {
	return fact.token
}

func (fact MintFact) Sender() base.Address {
	return fact.sender
}

func (fact MintFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.sender)

	return as, nil
}

func (fact MintFact) Items() []MintItem {
	return fact.items
}

func (fact MintFact) Currencies() []currency.CurrencyID {
	cs := make([]currency.CurrencyID, len(fact.items))

	for i := range fact.items {
		cs[i] = fact.items[i].Currency()
	}

	return cs
}

func (fact MintFact) Rebuild() MintFact {
	fact.h = fact.GenerateHash()

	return fact
}

type Mint struct {
	currency.BaseOperation
}

func NewMint(fact MintFact, fs []base.FactSign, memo string) (Mint, error) {
	bo, err := currency.NewBaseOperationFromFact(MintHint, fact, fs, memo)
	if err != nil {
		return Mint{}, err
	}
	return Mint{BaseOperation: bo}, nil
}
