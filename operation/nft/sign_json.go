package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type SignFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender mitumbase.Address `json:"sender"`
	Items  []SignItem        `json:"items"`
}

func (fact SignFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SignFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type SignFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *SignFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf SignFactJSONUnmarshaler

	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type signMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Sign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(signMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Sign) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
