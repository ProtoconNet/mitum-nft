package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	contract mitumbase.Address
	approved mitumbase.Address
	idx      uint64
	currency types.CurrencyID
}

func NewApproveItem(contract mitumbase.Address, approved mitumbase.Address, idx uint64, currency types.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		contract:   contract,
		approved:   approved,
		idx:        idx,
		currency:   currency,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.approved,
		it.currency,
	)
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.approved.Bytes(),
		util.Uint64ToBytes(it.idx),
		it.currency.Bytes(),
	)
}

func (it ApproveItem) Approved() mitumbase.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFT() uint64 {
	return it.idx
}

func (it ApproveItem) Currency() types.CurrencyID {
	return it.currency
}
