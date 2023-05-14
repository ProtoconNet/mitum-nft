package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
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
	contract   base.Address
	collection extensioncurrency.ContractID
	operator   base.Address
	mode       DelegateMode
	currency   currency.CurrencyID
}

func NewDelegateItem(contract base.Address, collection extensioncurrency.ContractID, operator base.Address, mode DelegateMode, currency currency.CurrencyID) DelegateItem {
	return DelegateItem{
		BaseHinter: hint.NewBaseHinter(DelegateItemHint),
		contract:   contract,
		collection: collection,
		operator:   operator,
		mode:       mode,
		currency:   currency,
	}
}

func (it DelegateItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.operator,
		it.mode,
		it.currency,
	)
}

func (it DelegateItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.operator.Bytes(),
		it.mode.Bytes(),
		it.currency.Bytes(),
	)
}

func (it DelegateItem) Contract() base.Address {
	return it.contract
}

func (it DelegateItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it DelegateItem) Operator() base.Address {
	return it.operator
}

func (it DelegateItem) Mode() DelegateMode {
	return it.mode
}

func (it DelegateItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.operator
	return as, nil
}

func (it DelegateItem) Currency() currency.CurrencyID {
	return it.currency
}
