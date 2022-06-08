package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	BurnFactType   = hint.Type("mitum-nft-burn-operation-fact")
	BurnFactHint   = hint.NewHint(BurnFactType, "v0.0.1")
	BurnFactHinter = BurnFact{BaseHinter: hint.NewBaseHinter(BurnFactHint)}
	BurnType       = hint.Type("mitum-nft-burn-operation")
	BurnHint       = hint.NewHint(BurnType, "v0.0.1")
	BurnHinter     = Burn{BaseOperation: operationHinter(BurnHint)}
)

var MaxBurnItems uint = 10

type BurnItem interface {
	hint.Hinter
	isvalid.IsValider
	NFTsItem
	Bytes() []byte
	Collection() extensioncurrency.ContractID
	Currency() currency.CurrencyID
	Rebuild() BurnItem
}

type BurnFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []BurnItem
}

func NewBurnFact(token []byte, sender base.Address, items []BurnItem) BurnFact {
	fact := BurnFact{
		BaseHinter: hint.NewBaseHinter(BurnFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact BurnFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact BurnFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact BurnFact) Bytes() []byte {
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

func (fact BurnFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for BurnFact")
	} else if n > int(MaxBurnItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxBurnItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nfts := fact.items[i].NFTs()

		for j := range nfts {
			if err := nfts[j].IsValid(nil); err != nil {
				return err
			}

			nft := nfts[j].String()
			if _, found := foundNFT[nft]; found {
				return isvalid.InvalidError.Errorf("duplicated nft found; %s", nft)
			}

			foundNFT[nft] = true
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact BurnFact) Token() []byte {
	return fact.token
}

func (fact BurnFact) Sender() base.Address {
	return fact.sender
}

func (fact BurnFact) Items() []BurnItem {
	return fact.items
}

func (fact BurnFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact BurnFact) Rebuild() BurnFact {
	items := make([]BurnItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Burn struct {
	currency.BaseOperation
}

func NewBurn(fact BurnFact, fs []base.FactSign, memo string) (Burn, error) {
	bo, err := currency.NewBaseOperationFromFact(BurnHint, fact, fs, memo)
	if err != nil {
		return Burn{}, err
	}

	return Burn{BaseOperation: bo}, nil
}
