package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	contract   base.Address
	collection extensioncurrency.ContractID
	approved   base.Address
	nft        nft.NFTID
	currency   currency.CurrencyID
}

func NewApproveItem(contract base.Address, collection extensioncurrency.ContractID, approved base.Address, n nft.NFTID, currency currency.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		contract:   contract,
		collection: collection,
		approved:   approved,
		nft:        n,
		currency:   currency,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.approved,
		it.nft,
		it.currency,
	)
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.approved.Bytes(),
		it.nft.Bytes(),
		it.currency.Bytes(),
	)
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFT() nft.NFTID {
	return it.nft
}

func (it ApproveItem) Currency() currency.CurrencyID {
	return it.currency
}
