package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DelegateItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  mitumbase.Address `json:"contract"`
	Delegatee mitumbase.Address `json:"delegatee"`
	Mode      DelegateMode      `json:"mode"`
	Currency  types.CurrencyID  `json:"currency"`
}

func (it DelegateItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegateItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Delegatee:  it.delegatee,
		Mode:       it.mode,
		Currency:   it.currency,
	})
}

type DelegateItemJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	Delegatee string    `json:"delegatee"`
	Mode      string    `json:"mode"`
	Currency  string    `json:"currency"`
}

func (it *DelegateItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of DelegateItem")

	var u DelegateItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Delegatee, u.Mode, u.Currency)
}
