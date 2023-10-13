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
	TransferFactHint = hint.MustNewHint("mitum-nft-transfer-operation-fact-v0.0.1")
	TransferHint     = hint.MustNewHint("mitum-nft-transfer-operation-v0.0.1")
)

var MaxTransferItems = 10

type TransferFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []TransferItem
}

func NewTransferFact(token []byte, sender mitumbase.Address, items []TransferItem) TransferFact {
	bf := mitumbase.NewBaseFact(TransferFactHint, token)

	fact := TransferFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for TransferFact")
	} else if l > int(MaxTransferItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxTransferItems)
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

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact TransferFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact TransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
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

func (fact TransferFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact TransferFact) Items() []TransferItem {
	return fact.items
}

func (fact TransferFact) Addresses() ([]mitumbase.Address, error) {
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

type Transfer struct {
	common.BaseOperation
}

func NewTransfer(fact TransferFact) (Transfer, error) {
	return Transfer{BaseOperation: common.NewBaseOperation(TransferHint, fact)}, nil
}
