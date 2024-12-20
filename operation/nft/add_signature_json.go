package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AddSignatureFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address       `json:"sender"`
	Items  []AddSignatureItem `json:"items"`
}

func (fact AddSignatureFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddSignatureFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type AddSignatureFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *AddSignatureFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf AddSignatureFactJSONUnmarshaler

	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type AddSignatureMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op AddSignature) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddSignatureMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AddSignature) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
