package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	ApproveAllAllow  = ApproveAllMode("allow")
	ApproveAllCancel = ApproveAllMode("cancel")
)

type ApproveAllMode string

func (mode ApproveAllMode) IsValid([]byte) error {
	if !(mode == ApproveAllAllow || mode == ApproveAllCancel) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("wrong approveAll mode, %v", mode))
	}

	return nil
}

func (mode ApproveAllMode) Bytes() []byte {
	return []byte(mode)
}

func (mode ApproveAllMode) String() string {
	return string(mode)
}

func (mode ApproveAllMode) Equal(cmode ApproveAllMode) bool {
	return string(mode) == string(cmode)
}

var ApproveAllItemHint = hint.MustNewHint("mitum-nft-approve-all-item-v0.0.1")

type ApproveAllItem struct {
	hint.BaseHinter
	contract base.Address
	approved base.Address
	mode     ApproveAllMode
	currency types.CurrencyID
}

func NewApproveAllItem(
	contract base.Address, approved base.Address, mode ApproveAllMode, currency types.CurrencyID) ApproveAllItem {
	return ApproveAllItem{
		BaseHinter: hint.NewBaseHinter(ApproveAllItemHint),
		contract:   contract,
		approved:   approved,
		mode:       mode,
		currency:   currency,
	}
}

func (it ApproveAllItem) IsValid([]byte) error {
	if it.approved.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("approved account %v is same with contract account", it.approved))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.approved,
		it.mode,
		it.currency,
	)
}

func (it ApproveAllItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.approved.Bytes(),
		it.mode.Bytes(),
		it.currency.Bytes(),
	)
}

func (it ApproveAllItem) Contract() base.Address {
	return it.contract
}

func (it ApproveAllItem) Approved() base.Address {
	return it.approved
}

func (it ApproveAllItem) Mode() ApproveAllMode {
	return it.mode
}

func (it ApproveAllItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveAllItem) Currency() types.CurrencyID {
	return it.currency
}
