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
	NFTTransferFactHint = hint.MustNewHint("mitum-nft-transfer-operation-fact-v0.0.1")
	NFTTransferHint     = hint.MustNewHint("mitum-nft-transfer-operation-v0.0.1")
)

var MaxNFTTransferItems = 10

type NFTTransferFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []NFTTransferItem
}

func NewNFTTransferFact(token []byte, sender mitumbase.Address, items []NFTTransferItem) NFTTransferFact {
	bf := mitumbase.NewBaseFact(NFTTransferFactHint, token)

	fact := NFTTransferFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact NFTTransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for NFTTransferFact")
	} else if l > int(MaxNFTTransferItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxNFTTransferItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return err
		}

		n := strconv.FormatUint(item.NFT(), 10)

		if _, found := founds[n]; found {
			return util.ErrInvalid.Errorf("duplicate nft found, %q", n)
		}

		founds[n] = struct{}{}
	}

	return nil
}

func (fact NFTTransferFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact NFTTransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact NFTTransferFact) Bytes() []byte {
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

func (fact NFTTransferFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact NFTTransferFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact NFTTransferFact) Items() []NFTTransferItem {
	return fact.items
}

func (fact NFTTransferFact) Addresses() ([]mitumbase.Address, error) {
	as := []mitumbase.Address{}

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.Sender())

	return as, nil
}

type NFTTransfer struct {
	common.BaseOperation
}

func NewNFTTransfer(fact NFTTransferFact) (NFTTransfer, error) {
	return NFTTransfer{BaseOperation: common.NewBaseOperation(NFTTransferHint, fact)}, nil
}

func (op *NFTTransfer) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
