package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *DelegateItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col string,
	ag string,
	md string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal DelegateItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	it.collection = currencybase.ContractID(col)
	it.mode = DelegateMode(md)
	it.currency = currencybase.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	operator, err := base.DecodeAddress(ag, enc)
	if err != nil {
		return e(err, "")
	}
	it.operator = operator

	return nil
}
