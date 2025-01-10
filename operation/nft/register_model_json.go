package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RegisterModelFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender    base.Address           `json:"sender"`
	Contract  base.Address           `json:"contract"`
	Name      types.CollectionName   `json:"name"`
	Royalty   types.PaymentParameter `json:"royalty"`
	URI       types.URI              `json:"uri"`
	Whitelist []base.Address         `json:"minter_whitelist"`
	Currency  ctypes.CurrencyID      `json:"currency"`
}

func (fact RegisterModelFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterModelFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Name:                  fact.name,
		Royalty:               fact.royalty,
		URI:                   fact.uri,
		Whitelist:             fact.minterWhitelist,
		Currency:              fact.currency,
	})
}

type RegisterModelFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender    string   `json:"sender"`
	Contract  string   `json:"contract"`
	Name      string   `json:"name"`
	Royalty   uint     `json:"royalty"`
	URI       string   `json:"uri"`
	Whitelist []string `json:"minter_whitelist"`
	Currency  string   `json:"currency"`
}

func (fact *RegisterModelFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u RegisterModelFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)
	if err := fact.unmarshal(enc, u.Sender, u.Contract, u.Name, u.Royalty, u.URI, u.Whitelist, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op RegisterModel) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *RegisterModel) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
