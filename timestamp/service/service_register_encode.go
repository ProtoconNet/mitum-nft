package service

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *ServiceRegisterFact) unmarshal(
	enc encoder.Encoder,
	sa,
	ta,
	svc,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal ServiceRegisterFact")

	fact.currency = currency.CurrencyID(cid)

	sender, err := base.DecodeAddress(sa, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender
	target, err := base.DecodeAddress(ta, enc)
	if err != nil {
		return e(err, "")
	}
	fact.target = target
	fact.service = extensioncurrency.ContractID(svc)

	return nil
}
