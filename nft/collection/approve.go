package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxApproveItems = 10

type ApproveItem interface {
	hint.Hinter
	isvalid.IsValider
	NFTsItem
	Bytes() []byte
	Approved() base.Address
	Addresses() ([]base.Address, error)
	Currency() currency.CurrencyID
	Rebuild() ApproveItem
}

var (
	ApproveFactType   = hint.Type("mitum-nft-approve-operation-fact")
	ApproveFactHint   = hint.NewHint(ApproveFactType, "v0.0.1")
	ApproveFactHinter = ApproveFact{BaseHinter: hint.NewBaseHinter(ApproveFactHint)}
	ApproveType       = hint.Type("mitum-nft-approve-operation")
	ApproveHint       = hint.NewHint(ApproveType, "v0.0.1")
	ApproveHinter     = Approve{BaseOperation: operationHinter(ApproveHint)}
)

type ApproveFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []ApproveItem
}

func NewApproveFact(token []byte, sender base.Address, items []ApproveItem) ApproveFact {
	fact := ApproveFact{
		BaseHinter: hint.NewBaseHinter(ApproveFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact ApproveFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact ApproveFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
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

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return isvalid.InvalidError.Errorf("empty items for ApproveFact")
	} else if l > int(MaxApproveItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", l, MaxApproveItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nfts := fact.items[i].NFTs()

		for j := range nfts {
			if err := nfts[j].IsValid(nil); err != nil {
				return err
			}

			n := nfts[j]
			if _, found := foundNFT[n]; found {
				return isvalid.InvalidError.Errorf("duplicated nft found; %q", n)
			}

			foundNFT[n] = true
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact ApproveFact) Token() []byte {
	return fact.token
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) NFTs() []nft.NFTID {
	ns := []nft.NFTID{}

	for i := range fact.items {
		ns = append(ns, fact.items[i].NFTs()...)
	}

	return ns
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	for i := range fact.items {
		as[i] = fact.items[i].Approved()
	}
	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact ApproveFact) Currencies() []currency.CurrencyID {
	cs := make([]currency.CurrencyID, len(fact.items))

	for i := range fact.items {
		cs[i] = fact.items[i].Currency()
	}

	return cs
}

func (fact ApproveFact) Rebuild() ApproveFact {
	items := make([]ApproveItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Approve struct {
	currency.BaseOperation
}

func NewApprove(fact ApproveFact, fs []base.FactSign, memo string) (Approve, error) {
	bo, err := currency.NewBaseOperationFromFact(ApproveHint, fact, fs, memo)
	if err != nil {
		return Approve{}, err
	}
	return Approve{BaseOperation: bo}, nil
}
