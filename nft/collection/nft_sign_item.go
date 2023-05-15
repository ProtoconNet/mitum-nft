package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	contract   base.Address
	collection extensioncurrency.ContractID
	nft        uint64
	currency   currency.CurrencyID
}

func NewNFTSignItem(contract base.Address, collection extensioncurrency.ContractID, n uint64, currency currency.CurrencyID) NFTSignItem {
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

func (it NFTSignItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it NFTSignItem) Collection() extensioncurrency.ContractID {
	return it.collection
}
