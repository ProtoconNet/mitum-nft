package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

// type CollectionRegisterFormJSONMarshaler struct {
// 	hint.BaseHinter
// 	Contract   base.Address                 `json:"contract"`
// 	Collection extensioncurrency.ContractID `json:"collection"`
// 	Name       CollectionName               `json:"name"`
// 	Royalty    nft.PaymentParameter         `json:"royalty"`
// 	URI        nft.URI                      `json:"uri"`
// 	Whites     []base.Address               `json:"whites"`
// }

// func (form CollectionRegisterForm) MarshalJSON() ([]byte, error) {
// 	return util.MarshalJSON(CollectionRegisterFormJSONMarshaler{
// 		BaseHinter: form.BaseHinter,
// 		Contract:   form.contract,
// 		Collection: form.collection,
// 		Name:       form.name,
// 		Royalty:    form.royalty,
// 		URI:        form.uri,
// 		Whites:     form.whites,
// 	})
// }

// type CollectionRegisterFormJSONUnmarshaler struct {
// 	Hint       hint.Hint `json:"_hint"`
// 	Contract   string    `json:"contract"`
// 	Collection string    `json:"collection"`
// 	Name       string    `json:"name"`
// 	Royalty    uint      `json:"royalty"`
// 	URI        string    `json:"uri"`
// 	Whites     []string  `json:"whites"`
// }

// func (form *CollectionRegisterForm) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
// 	e := util.StringErrorFunc("failed to decode json of CollectionRegisterForm")

// 	var u CollectionRegisterFormJSONUnmarshaler
// 	if err := enc.Unmarshal(b, &u); err != nil {
// 		return e(err, "")
// 	}

// 	return form.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.Name, u.Royalty, u.URI, u.Whites)
// }

type CollectionRegisterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender     base.Address                 `json:"sender"`
	Contract   base.Address                 `json:"contract"`
	Collection extensioncurrency.ContractID `json:"collection"`
	Name       CollectionName               `json:"name"`
	Royalty    nft.PaymentParameter         `json:"royalty"`
	URI        nft.URI                      `json:"uri"`
	Whites     []base.Address               `json:"whites"`
	Currency   currency.CurrencyID          `json:"currency"`
}

func (fact CollectionRegisterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionRegisterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Collection:            fact.collection,
		Name:                  fact.name,
		Royalty:               fact.royalty,
		URI:                   fact.uri,
		Whites:                fact.whitelist,
		Currency:              fact.currency,
	})
}

type CollectionRegisterFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender     string   `json:"sender"`
	Contract   string   `json:"contract"`
	Collection string   `json:"collection"`
	Name       string   `json:"name"`
	Royalty    uint     `json:"royalty"`
	URI        string   `json:"uri"`
	Whites     []string `json:"whites"`
	Currency   string   `json:"currency"`
}

func (fact *CollectionRegisterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionRegisterFact")

	var u CollectionRegisterFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Contract, u.Collection, u.Name, u.Royalty, u.URI, u.Whites, u.Currency)
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
