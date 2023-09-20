package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *DelegateItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	cAdr, dAdr, md, cid string,
) error {
	e := util.StringError("failed to unmarshal DelegateItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	it.mode = DelegateMode(md)
	it.currency = types.CurrencyID(cid)

	switch a, err := mitumbase.DecodeAddress(cAdr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	delegatee, err := mitumbase.DecodeAddress(dAdr, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.delegatee = delegatee

	return nil
}
