package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type UpdateCollectionPolicyFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender    mitumbase.Address        `json:"sender"`
	Contract  mitumbase.Address        `json:"contract"`
	Name      types.CollectionName     `json:"name"`
	Royalty   types.PaymentParameter   `json:"royalty"`
	URI       types.URI                `json:"uri"`
	Whitelist []mitumbase.Address      `json:"whitelist"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (fact UpdateCollectionPolicyFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(UpdateCollectionPolicyFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Name:                  fact.name,
		Royalty:               fact.royalty,
		URI:                   fact.uri,
		Whitelist:             fact.whitelist,
		Currency:              fact.currency,
	})
}

type UpdateCollectionPolicyFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender    string   `json:"sender"`
	Contract  string   `json:"contract"`
	Name      string   `json:"name"`
	Royalty   uint     `json:"royalty"`
	URI       string   `json:"uri"`
	Whitelist []string `json:"whitelist"`
	Currency  string   `json:"currency"`
}

func (fact *UpdateCollectionPolicyFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of UpdateCollectionPolicyFact")

	var u UpdateCollectionPolicyFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Contract, u.Name, u.Royalty, u.URI, u.Whitelist, u.Currency)
}

type updateCollectionPolicyMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op UpdateCollectionPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(updateCollectionPolicyMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *UpdateCollectionPolicy) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of UpdateCollectionPolicy")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
