package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	contract base.Address
	approved base.Address
	nftIdx   uint64
	currency types.CurrencyID
}

func NewApproveItem(
	contract base.Address, approved base.Address, nftIdx uint64, currency types.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		contract:   contract,
		approved:   approved,
		nftIdx:     nftIdx,
		currency:   currency,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	if it.approved.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("approved %v is same with contract contract", it.approved))
	}

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
		util.Uint64ToBytes(it.nftIdx),
		it.currency.Bytes(),
	)
}

func (it ApproveItem) Contract() base.Address {
	return it.contract
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFTIdx() uint64 {
	return it.nftIdx
}

func (it ApproveItem) Currency() types.CurrencyID {
	return it.currency
}
