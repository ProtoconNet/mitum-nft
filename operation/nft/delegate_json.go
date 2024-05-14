package nft

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DelegateFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender mitumbase.Address `json:"sender"`
	Items  []DelegateItem    `json:"items"`
}

func (fact DelegateFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegateFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type DelegateFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *DelegateFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u DelegateFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unmarshal(enc, u.Sender, u.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type delegateMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Delegate) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(delegateMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Delegate) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
