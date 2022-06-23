package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	DelegateAllow  = DelegateMode("allow")
	DelegateCancel = DelegateMode("cancel")
)

type DelegateMode string

func (mode DelegateMode) Bytes() []byte {
	return []byte(mode)
}

func (mode DelegateMode) String() string {
	return string(mode)
}

func (mode DelegateMode) Equal(cmode DelegateMode) bool {
	return string(mode) == string(cmode)
}

func (mode DelegateMode) IsValid([]byte) error {
	if !(mode == DelegateAllow || mode == DelegateCancel) {
		return isvalid.InvalidError.Errorf("wrong delegate mode; %s", mode)
	}

	return nil
}

var (
	DelegateItemType   = hint.Type("mitum-nft-delegate-item")
	DelegateItemHint   = hint.NewHint(DelegateItemType, "v0.0.1")
	DelegateItemHinter = DelegateItem{BaseHinter: hint.NewBaseHinter(DelegateItemHint)}
)

type DelegateItem struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	agent      base.Address
	mode       DelegateMode
	cid        currency.CurrencyID
}

func NewDelegateItem(symbol extensioncurrency.ContractID, agent base.Address, mode DelegateMode, cid currency.CurrencyID) DelegateItem {
	return DelegateItem{
		BaseHinter: hint.NewBaseHinter(DelegateItemHint),
		collection: symbol,
		agent:      agent,
		mode:       mode,
		cid:        cid,
	}
}

func (it DelegateItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.collection.Bytes(),
		it.agent.Bytes(),
		it.mode.Bytes(),
		it.cid.Bytes(),
	)
}

func (it DelegateItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		it.BaseHinter,
		it.collection,
		it.agent,
		it.mode,
		it.cid); err != nil {
		return err
	}

	return nil
}

func (it DelegateItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it DelegateItem) Agent() base.Address {
	return it.agent
}

func (it DelegateItem) Mode() DelegateMode {
	return it.mode
}

func (it DelegateItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.agent
	return as, nil
}

func (it DelegateItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it DelegateItem) Rebuild() DelegateItem {
	return it
}
