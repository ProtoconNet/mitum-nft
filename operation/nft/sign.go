package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	SignFactHint = hint.MustNewHint("mitum-nft-sign-operation-fact-v0.0.1")
	SignHint     = hint.MustNewHint("mitum-nft-sign-operation-v0.0.1")
)

var MaxSignItems = 100

type SignFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []SignItem
}

func NewSignFact(token []byte, sender mitumbase.Address, items []SignItem) SignFact {
	bf := mitumbase.NewBaseFact(SignFactHint, token)
	fact := SignFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact SignFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for SignFact")))
	} else if l > int(MaxSignItems) {
		return common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxSignItems))
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		nid := strconv.FormatUint(item.NFT(), 10)
		if _, found := founds[nid]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf("nft idx %v in contract account %v", item.NFT(), item.contract)))
		}

		founds[nid] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact SignFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact SignFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SignFact) Bytes() []byte {
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

func (fact SignFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact SignFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact SignFact) Items() []SignItem {
	return fact.items
}

func (fact SignFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type Sign struct {
	common.BaseOperation
}

func NewSign(fact SignFact) (Sign, error) {
	return Sign{BaseOperation: common.NewBaseOperation(SignHint, fact)}, nil
}
