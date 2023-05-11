package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CollectionPolicyUpdaterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender     base.Address                 `json:"sender"`
	Collection extensioncurrency.ContractID `json:"collection"`
	Name       CollectionName               `json:"name"`
	Royalty    nft.PaymentParameter         `json:"royalty"`
	URI        nft.URI                      `json:"uri"`
	Whitelist  []base.Address               `json:"whitelist"`
	Currency   currency.CurrencyID          `json:"currency"`
}

func (fact CollectionPolicyUpdaterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionPolicyUpdaterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Collection:            fact.collection,
		Name:                  fact.name,
		Royalty:               fact.royalty,
		URI:                   fact.uri,
		Whitelist:             fact.whitelist,
		Currency:              fact.currency,
	})
}

type CollectionPolicyUpdaterFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender     string   `json:"sender"`
	Collection string   `json:"collection"`
	Name       string   `json:"name"`
	Royalty    uint     `json:"royalty"`
	URI        string   `json:"uri"`
	Whitelist  []string `json:"whitelist"`
	Currency   string   `json:"currency"`
}

func (fact *CollectionPolicyUpdaterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionPolicyUpdaterFact")

	var u CollectionPolicyUpdaterFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Collection, u.Name, u.Royalty, u.URI, u.Whitelist, u.Currency)
}

type collectionPolicyUpdaterMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op CollectionPolicyUpdater) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(collectionPolicyUpdaterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CollectionPolicyUpdater) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionPolicyUpdater")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
