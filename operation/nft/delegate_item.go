package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
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
	contract   mitumbase.Address
	collection types.ContractID
	operator   mitumbase.Address
	mode       DelegateMode
	currency   types.CurrencyID
}

func NewDelegateItem(contract mitumbase.Address, collection types.ContractID, operator mitumbase.Address, mode DelegateMode, currency types.CurrencyID) DelegateItem {
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

func (it DelegateItem) Contract() mitumbase.Address {
	return it.contract
}

func (it DelegateItem) Collection() types.ContractID {
	return it.collection
}

func (it DelegateItem) Operator() mitumbase.Address {
	return it.operator
}

func (it DelegateItem) Mode() DelegateMode {
	return it.mode
}

func (it DelegateItem) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = it.operator
	return as, nil
}

func (it DelegateItem) Currency() types.CurrencyID {
	return it.currency
}
