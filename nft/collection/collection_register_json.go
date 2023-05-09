package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CollectionRegisterFormJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address                 `json:"contract"`
	Symbol   extensioncurrency.ContractID `json:"symbol"`
	Name     CollectionName               `json:"name"`
	Royalty  nft.PaymentParameter         `json:"royalty"`
	URI      nft.URI                      `json:"uri"`
	Whites   []base.Address               `json:"whites"`
}

func (form CollectionRegisterForm) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionRegisterFormJSONMarshaler{
		BaseHinter: form.BaseHinter,
		Contract:   form.contract,
		Symbol:     form.symbol,
		Name:       form.name,
		Royalty:    form.royalty,
		URI:        form.uri,
		Whites:     form.whites,
	})
}

type CollectionRegisterFormJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Symbol   string    `json:"symbol"`
	Name     string    `json:"name"`
	Royalty  uint      `json:"royalty"`
	URI      string    `json:"uri"`
	Whites   []string  `json:"whites"`
}

func (form *CollectionRegisterForm) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionRegisterForm")

	var u CollectionRegisterFormJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return form.unmarshal(enc, u.Hint, u.Contract, u.Symbol, u.Name, u.Royalty, u.URI, u.Whites)
}

type CollectionRegisterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address           `json:"sender"`
	Form     CollectionRegisterForm `json:"form"`
	Currency currency.CurrencyID    `json:"currency"`
}

func (fact CollectionRegisterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionRegisterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Form:                  fact.form,
		Currency:              fact.currency,
	})
}

type CollectionRegisterFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string          `json:"sender"`
	Form     json.RawMessage `json:"form"`
	Currency string          `json:"currency"`
}

func (fact *CollectionRegisterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionRegisterFact")

	var u CollectionRegisterFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Form, u.Currency)
}

type collectionRegisterMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op CollectionRegister) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(collectionRegisterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CollectionRegister) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrecyRegister")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
