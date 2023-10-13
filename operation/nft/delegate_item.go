package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DelegateAllow  = DelegateMode("allow")
	DelegateCancel = DelegateMode("cancel")
)

type DelegateMode string

func (mode DelegateMode) IsValid([]byte) error {
	if !(mode == DelegateAllow || mode == DelegateCancel) {
		return util.ErrInvalid.Errorf("wrong delegate mode, %q", mode)
	}

	return nil
}

func (mode DelegateMode) Bytes() []byte {
	return []byte(mode)
}

func (mode DelegateMode) String() string {
	return string(mode)
}

func (mode DelegateMode) Equal(cmode DelegateMode) bool {
	return string(mode) == string(cmode)
}

var DelegateItemHint = hint.MustNewHint("mitum-nft-delegate-item-v0.0.1")

type DelegateItem struct {
	hint.BaseHinter
	contract  mitumbase.Address
	delegatee mitumbase.Address
	mode      DelegateMode
	currency  types.CurrencyID
}

func NewDelegateItem(contract mitumbase.Address, operator mitumbase.Address, mode DelegateMode, currency types.CurrencyID) DelegateItem {
	return DelegateItem{
		BaseHinter: hint.NewBaseHinter(DelegateItemHint),
		contract:   contract,
		delegatee:  operator,
		mode:       mode,
		currency:   currency,
	}
}

func (it DelegateItem) IsValid([]byte) error {
	if it.contract == nil {
		return errors.Errorf("contract is nil in DelegateItem")
	}
	if it.delegatee == nil {
		return errors.Errorf("delegatee is nil in DelegateItem")
	}
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.delegatee,
		it.mode,
		it.currency,
	)
}

func (it DelegateItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.delegatee.Bytes(),
		it.mode.Bytes(),
		it.currency.Bytes(),
	)
}

func (it DelegateItem) Contract() mitumbase.Address {
	return it.contract
}

func (it DelegateItem) Delegatee() mitumbase.Address {
	return it.delegatee
}

func (it DelegateItem) Mode() DelegateMode {
	return it.mode
}

func (it DelegateItem) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = it.delegatee
	return as, nil
}

func (it DelegateItem) Currency() types.CurrencyID {
	return it.currency
}
