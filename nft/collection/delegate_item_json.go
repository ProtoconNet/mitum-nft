package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DelegateItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address                 `json:"contract"`
	Collection extensioncurrency.ContractID `json:"collection"`
	Operator   base.Address                 `json:"operator"`
	Mode       DelegateMode                 `json:"mode"`
	Currency   currency.CurrencyID          `json:"currency"`
}

func (it DelegateItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegateItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Collection: it.collection,
		Operator:   it.operator,
		Mode:       it.mode,
		Currency:   it.currency,
	})
}

type DelegateItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Collection string    `json:"collection"`
	Operator   string    `json:"operator"`
	Mode       string    `json:"mode"`
	Currency   string    `json:"currency"`
}

func (it *DelegateItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of DelegateItem")

	var u DelegateItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.Operator, u.Mode, u.Currency)
}
