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
	AddSignatureFactHint = hint.MustNewHint("mitum-nft-add-signature-operation-fact-v0.0.1")
	AddSignatureHint     = hint.MustNewHint("mitum-nft-add-signature-operation-v0.0.1")
)

var MaxAddSignatureItems = 100

type AddSignatureFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []AddSignatureItem
}

func NewAddSignatureFact(token []byte, sender mitumbase.Address, items []AddSignatureItem) AddSignatureFact {
	bf := mitumbase.NewBaseFact(AddSignatureFactHint, token)
	fact := AddSignatureFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AddSignatureFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for AddSignatureFact")))
	} else if l > int(MaxAddSignatureItems) {
		return common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxAddSignatureItems))
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

func (fact AddSignatureFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AddSignatureFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddSignatureFact) Bytes() []byte {
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

func (fact AddSignatureFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact AddSignatureFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact AddSignatureFact) Items() []AddSignatureItem {
	return fact.items
}

func (fact AddSignatureFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type AddSignature struct {
	common.BaseOperation
}

func NewSign(fact AddSignatureFact) (AddSignature, error) {
	return AddSignature{BaseOperation: common.NewBaseOperation(AddSignatureHint, fact)}, nil
}
