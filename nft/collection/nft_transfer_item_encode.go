package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *NFTTransferItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col,
	rc string,
	nid uint64,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTTransferItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.collection = extensioncurrency.ContractID(col)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	receiver, err := base.DecodeAddress(rc, enc)
	if err != nil {
		return e(err, "")
	}
	it.receiver = receiver
	it.nft = nid
	it.currency = currency.CurrencyID(cid)

	return nil
}
