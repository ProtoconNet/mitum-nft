package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTTransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type NFTTransferItem struct {
	hint.BaseHinter
	contract   base.Address
	collection currencybase.ContractID
	receiver   base.Address
	nft        uint64
	currency   currencybase.CurrencyID
}

func NewNFTTransferItem(contract base.Address, collection currencybase.ContractID, receiver base.Address, nft uint64, currency currencybase.CurrencyID) NFTTransferItem {
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

func (it NFTTransferItem) Contract() base.Address {
	return it.contract
}

func (it NFTTransferItem) Collection() currencybase.ContractID {
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

func (it NFTTransferItem) NFT() nft.NFTID {
	nftID := nft.NFTID(it.nft)
	return nftID
}

func (it NFTTransferItem) Currency() currencybase.CurrencyID {
	return it.currency
}
