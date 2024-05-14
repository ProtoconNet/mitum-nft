package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *SignItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca string,
	nft uint64,
	cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = types.CurrencyID(cid)
	switch a, err := mitumbase.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	it.nft = nft

	return nil
}
