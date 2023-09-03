package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	SignFactHint = hint.MustNewHint("mitum-nft-sign-operation-fact-v0.0.1")
	SignHint     = hint.MustNewHint("mitum-nft-sign-operation-v0.0.1")
)

var MaxSignItems = 10

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
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for SignFact")
	} else if l > int(MaxSignItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxSignItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return err
		}

		nid := strconv.FormatUint(item.NFT(), 10)
		if _, found := founds[nid]; found {
			return util.ErrInvalid.Errorf("duplicate nft found, %q", item.NFT())
		}

		founds[nid] = struct{}{}
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

func (op *Sign) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
