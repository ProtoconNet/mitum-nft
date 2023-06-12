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
	NFTSignFactHint = hint.MustNewHint("mitum-nft-sign-operation-fact-v0.0.1")
	NFTSignHint     = hint.MustNewHint("mitum-nft-sign-operation-v0.0.1")
)

var MaxNFTSignItems = 10

type NFTSignFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []NFTSignItem
}

func NewNFTSignFact(token []byte, sender mitumbase.Address, items []NFTSignItem) NFTSignFact {
	bf := mitumbase.NewBaseFact(NFTSignFactHint, token)
	fact := NFTSignFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact NFTSignFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for NFTSignFact")
	} else if l > int(MaxNFTSignItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxNFTSignItems)
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

func (fact NFTSignFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact NFTSignFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact NFTSignFact) Bytes() []byte {
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

func (fact NFTSignFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact NFTSignFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact NFTSignFact) Items() []NFTSignItem {
	return fact.items
}

func (fact NFTSignFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type NFTSign struct {
	common.BaseOperation
}

func NewNFTSign(fact NFTSignFact) (NFTSign, error) {
	return NFTSign{BaseOperation: common.NewBaseOperation(NFTSignHint, fact)}, nil
}

func (op *NFTSign) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
