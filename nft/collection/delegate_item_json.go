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
	Collection extensioncurrency.ContractID `json:"collection"`
	Agent      base.Address                 `json:"agent"`
	Mode       DelegateMode                 `json:"mode"`
	Currency   currency.CurrencyID          `json:"currency"`
}

func (it DelegateItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegateItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Collection: it.collection,
		Agent:      it.operator,
		Mode:       it.mode,
		Currency:   it.currency,
	})
}

type DelegateItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Collection string    `json:"collection"`
	Agent      string    `json:"agent"`
	Mode       string    `json:"mode"`
	Currency   string    `json:"currency"`
}

func (it *DelegateItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of DelegateItem")

	var u DelegateItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, u.Hint, u.Collection, u.Agent, u.Mode, u.Currency)
}
