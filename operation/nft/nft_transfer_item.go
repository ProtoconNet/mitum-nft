package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTTransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type NFTTransferItem struct {
	hint.BaseHinter
	contract   mitumbase.Address
	collection types.ContractID
	receiver   mitumbase.Address
	nft        uint64
	currency   types.CurrencyID
}

func NewNFTTransferItem(contract mitumbase.Address, collection types.ContractID, receiver mitumbase.Address, nft uint64, currency types.CurrencyID) NFTTransferItem {
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
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.receiver.Bytes(),
		util.Uint64ToBytes(it.nft),
		it.currency.Bytes(),
	)
}

func (it NFTTransferItem) Contract() mitumbase.Address {
	return it.contract
}

func (it NFTTransferItem) Collection() types.ContractID {
	return it.collection
}

func (it NFTTransferItem) Receiver() mitumbase.Address {
	return it.receiver
}

func (it NFTTransferItem) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it NFTTransferItem) NFT() uint64 {
	return it.nft
}

func (it NFTTransferItem) Currency() types.CurrencyID {
	return it.currency
}
