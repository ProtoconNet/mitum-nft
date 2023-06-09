package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	contract   base.Address
	collection currencybase.ContractID
	nft        uint64
	currency   currencybase.CurrencyID
}

func NewNFTSignItem(contract base.Address, collection currencybase.ContractID, n uint64, currency currencybase.CurrencyID) NFTSignItem {
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

func (it NFTSignItem) NFT() nft.NFTID {
	return nft.NFTID(it.nft)
}

func (it NFTSignItem) Contract() base.Address {
	return it.contract
}

func (it NFTSignItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it NFTSignItem) Collection() currencybase.ContractID {
	return it.collection
}
