package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *DelegateItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	ag string,
	md string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal DelegateItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	it.collection = extensioncurrency.ContractID(col)
	it.mode = DelegateMode(md)
	it.currency = currency.CurrencyID(cid)

	operator, err := base.DecodeAddress(ag, enc)
	if err != nil {
		return e(err, "")
	}
	it.operator = operator

	return nil
}
