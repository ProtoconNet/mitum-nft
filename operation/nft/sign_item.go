package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var SignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type SignItem struct {
	hint.BaseHinter
	contract mitumbase.Address
	nft      uint64
	currency types.CurrencyID
}

func NewSignItem(contract mitumbase.Address, n uint64, currency types.CurrencyID) SignItem {
	return SignItem{
		BaseHinter: hint.NewBaseHinter(SignItemHint),
		contract:   contract,
		nft:        n,
		currency:   currency,
	}
}

func (it SignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		util.Uint64ToBytes(it.nft),
		it.currency.Bytes(),
	)
}

func (it SignItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.currency)
}

func (it SignItem) NFT() uint64 {
	return it.nft
}

func (it SignItem) Contract() mitumbase.Address {
	return it.contract
}

func (it SignItem) Currency() types.CurrencyID {
	return it.currency
}
