package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	contract   mitumbase.Address
	collection types.ContractID
	nft        uint64
	currency   types.CurrencyID
}

func NewNFTSignItem(contract mitumbase.Address, collection types.ContractID, n uint64, currency types.CurrencyID) NFTSignItem {
	return NFTSignItem{
		BaseHinter: hint.NewBaseHinter(NFTSignItemHint),
		contract:   contract,
		collection: collection,
		nft:        n,
		currency:   currency,
	}
}

func (it NFTSignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		util.Uint64ToBytes(it.nft),
		it.currency.Bytes(),
	)
}

func (it NFTSignItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.collection, it.currency)
}

func (it NFTSignItem) NFT() uint64 {
	return it.nft
}

func (it NFTSignItem) Contract() mitumbase.Address {
	return it.contract
}

func (it NFTSignItem) Currency() types.CurrencyID {
	return it.currency
}

func (it NFTSignItem) Collection() types.ContractID {
	return it.collection
}
