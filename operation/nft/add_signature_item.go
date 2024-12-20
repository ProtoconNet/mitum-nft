package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AddSignatureItemHint = hint.MustNewHint("mitum-nft-add-signature-item-v0.0.1")

type AddSignatureItem struct {
	hint.BaseHinter
	contract base.Address
	nftIdx   uint64
	currency types.CurrencyID
}

func NewAddSignatureItem(contract base.Address, nfxIdx uint64, currency types.CurrencyID) AddSignatureItem {
	return AddSignatureItem{
		BaseHinter: hint.NewBaseHinter(AddSignatureItemHint),
		contract:   contract,
		nftIdx:     nfxIdx,
		currency:   currency,
	}
}

func (it AddSignatureItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		util.Uint64ToBytes(it.nftIdx),
		it.currency.Bytes(),
	)
}

func (it AddSignatureItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.currency,
	)
}

func (it AddSignatureItem) NFT() uint64 {
	return it.nftIdx
}

func (it AddSignatureItem) Contract() base.Address {
	return it.contract
}

func (it AddSignatureItem) Currency() types.CurrencyID {
	return it.currency
}
