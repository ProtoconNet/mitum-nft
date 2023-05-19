package service

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AppendFact) unmarshal(
	enc encoder.Encoder,
	sa string,
	ta string,
	svc,
	pid string,
	rqts uint64,
	data string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal AppendFact")

	switch sender, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = sender
	}

	switch target, err := base.DecodeAddress(ta, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.target = target
	}

	fact.service = extensioncurrency.ContractID(svc)
	fact.projectID = pid
	fact.requestTimeStamp = rqts
	fact.data = data
	fact.currency = currency.CurrencyID(cid)

	return nil
}
