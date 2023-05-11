package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *NFTSignItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col string,
	nft uint64,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTSignItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = currency.CurrencyID(cid)
	it.collection = extensioncurrency.ContractID(col)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	it.nft = nft

	return nil
}
