package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxMintItems = 100

var (
	MintFactHint = hint.MustNewHint("mitum-nft-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-nft-mint-operation-v0.0.1")
)

type MintFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []MintItem
}

func NewMintFact(token []byte, sender mitumbase.Address, items []MintItem) MintFact {
	bf := mitumbase.NewBaseFact(MintFactHint, token)
	fact := MintFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid(b []byte) error {
	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrArrayLen.Wrap(errors.Errorf("empty items for MintFact"))
	} else if l > int(MaxMintItems) {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxMintItems)))
	}

	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact MintFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MintFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))

	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact MintFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact MintFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact MintFact) Addresses() ([]mitumbase.Address, error) {
	as := []mitumbase.Address{}

	for _, item := range fact.items {
		if ads, err := item.Addresses(); err != nil {
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

type Mint struct {
	common.BaseOperation
}

func NewMint(fact MintFact) (Mint, error) {
	return Mint{BaseOperation: common.NewBaseOperation(MintHint, fact)}, nil
}
