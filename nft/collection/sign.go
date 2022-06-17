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

var (
	SignFactType   = hint.Type("mitum-nft-sign-operation-fact")
	SignFactHint   = hint.NewHint(SignFactType, "v0.0.1")
	SignFactHinter = SignFact{BaseHinter: hint.NewBaseHinter(SignFactHint)}
	SignType       = hint.Type("mitum-nft-sign-operation")
	SignHint       = hint.NewHint(SignType, "v0.0.1")
	SignHinter     = Sign{BaseOperation: operationHinter(SignHint)}
)

var MaxSignItems uint = 10

type SignFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []SignItem
}

func NewSignFact(token []byte, sender base.Address, items []SignItem) SignFact {
	fact := SignFact{
		BaseHinter: hint.NewBaseHinter(SignFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact SignFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact SignFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SignFact) Bytes() []byte {
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

func (fact SignFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return isvalid.InvalidError.Errorf("empty items for SignFact")
	} else if l > int(MaxSignItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", l, MaxSignItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		n := fact.items[i].NFT()
		if err := n.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundNFT[n]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found; %q", n)
		}

		foundNFT[n] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact SignFact) Token() []byte {
	return fact.token
}

func (fact SignFact) Sender() base.Address {
	return fact.sender
}

func (fact SignFact) Items() []SignItem {
	return fact.items
}

func (fact SignFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact SignFact) Rebuild() SignFact {
	items := make([]SignItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Sign struct {
	currency.BaseOperation
}

func NewSign(fact SignFact, fs []base.FactSign, memo string) (Sign, error) {
	bo, err := currency.NewBaseOperationFromFact(SignHint, fact, fs, memo)
	if err != nil {
		return Sign{}, err
	}

	return Sign{BaseOperation: bo}, nil
}
