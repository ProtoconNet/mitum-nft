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
	TransferFactType   = hint.Type("mitum-nft-tranfer-operation-fact")
	TransferFactHint   = hint.NewHint(TransferFactType, "v0.0.1")
	TransferFactHinter = TransferFact{BaseHinter: hint.NewBaseHinter(TransferFactHint)}
	TransferType       = hint.Type("mitum-nft-transfer-operation")
	TransferHint       = hint.NewHint(TransferType, "v0.0.1")
	TransferHinter     = Transfer{BaseOperation: operationHinter(TransferHint)}
)

var MaxTransferItems uint = 10

type NFTsItem interface {
	NFTs() []nft.NFTID
}

type TransferItem interface {
	hint.Hinter
	isvalid.IsValider
	NFTsItem
	Bytes() []byte
	Receiver() base.Address
	Currency() currency.CurrencyID
	Addresses() []base.Address
	Rebuild() TransferItem
}

type TransferFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []TransferItem
}

func NewTransferFact(token []byte, sender base.Address, items []TransferItem) TransferFact {
	fact := TransferFact{
		BaseHinter: hint.NewBaseHinter(TransferFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact TransferFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact TransferFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
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

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for TransferFact")
	} else if n > int(MaxTransferItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxTransferItems)
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

func (fact TransferFact) Token() []byte {
	return fact.token
}

func (fact TransferFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFact) Items() []TransferItem {
	return fact.items
}

func (fact TransferFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	for i := range fact.items {
		as = append(as, fact.items[i].Addresses()...)
	}

	as = append(as, fact.Sender())

	return as, nil
}

func (fact TransferFact) Rebuild() TransferFact {
	items := make([]TransferItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Transfer struct {
	currency.BaseOperation
}

func NewTransfer(fact TransferFact, fs []base.FactSign, memo string) (Transfer, error) {
	bo, err := currency.NewBaseOperationFromFact(TransferHint, fact, fs, memo)
	if err != nil {
		return Transfer{}, err
	}

	return Transfer{BaseOperation: bo}, nil
}
