package service

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type ServiceRegisterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address                 `json:"sender"`
	Target   base.Address                 `json:"target"`
	Service  extensioncurrency.ContractID `json:"service"`
	Currency currency.CurrencyID          `json:"currency"`
}

func (fact ServiceRegisterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ServiceRegisterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Target:                fact.target,
		Service:               fact.service,
		Currency:              fact.currency,
	})
}

type ServiceRegisterFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Target   string `json:"target"`
	Service  string `json:"service"`
	Currency string `json:"currency"`
}

func (fact *ServiceRegisterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ServiceRegisterFact")

	var u ServiceRegisterFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Target, u.Service, u.Currency)
}

type serviceRegisterMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op ServiceRegister) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(serviceRegisterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *ServiceRegister) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ServiceRegister")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
