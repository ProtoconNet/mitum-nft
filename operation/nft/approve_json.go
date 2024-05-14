package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type ApproveFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender mitumbase.Address `json:"sender"`
	Items  []ApproveItem     `json:"items"`
}

func (fact ApproveFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type ApproveFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *ApproveFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf ApproveFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unmarshal(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type approveMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Approve) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(approveMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Approve) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
