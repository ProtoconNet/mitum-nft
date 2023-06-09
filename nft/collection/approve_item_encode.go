package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *ApproveItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col,
	ap string,
	idx uint64,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal ApproveItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = currencybase.CurrencyID(cid)
	it.collection = currencybase.ContractID(col)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return e(err, "")
	}
	it.approved = approved
	it.idx = idx

	return nil
}
