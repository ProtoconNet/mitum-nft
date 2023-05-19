package collection

import (
	"fmt"
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTTransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type NFTTransferItem struct {
	hint.BaseHinter
	contract   base.Address
	collection extensioncurrency.ContractID
	receiver   base.Address
	nft        uint64
	currency   currency.CurrencyID
}

func NewNFTTransferItem(contract base.Address, collection extensioncurrency.ContractID, receiver base.Address, nft uint64, currency currency.CurrencyID) NFTTransferItem {
	return NFTTransferItem{
		BaseHinter: hint.NewBaseHinter(NFTTransferItemHint),
		contract:   contract,
		collection: collection,
		receiver:   receiver,
		nft:        nft,
		currency:   currency,
	}
}

func (it NFTTransferItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.receiver, it.currency)
}

func (it NFTTransferItem) Bytes() []byte {
	fmt.Printf("%08b", util.ConcatBytesSlice(
		util.Uint64ToBytes(it.nft),
	))
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.receiver.Bytes(),
		util.Uint64ToBytes(it.nft),
		it.currency.Bytes(),
	)
}

func (it NFTTransferItem) Contract() base.Address {
	return it.contract
}

func (it NFTTransferItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it NFTTransferItem) Receiver() base.Address {
	return it.receiver
}

func (it NFTTransferItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it NFTTransferItem) NFT() uint64 {
	return it.nft
}

func (it NFTTransferItem) Currency() currency.CurrencyID {
	return it.currency
}
