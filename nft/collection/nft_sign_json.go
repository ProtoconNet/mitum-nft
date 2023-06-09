package collection

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type NFTSignFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address  `json:"sender"`
	Items  []NFTSignItem `json:"items"`
}

func (fact NFTSignFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTSignFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type NFTSignFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *NFTSignFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTSignFact")

	var uf NFTSignFactJSONUnmarshaler

	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, uf.Sender, uf.Items)
}

type nftSignMarshaler struct {
	currencybase.BaseOperationJSONMarshaler
}

func (op NFTSign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(nftSignMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *NFTSign) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTSign")

	var ubo currencybase.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
