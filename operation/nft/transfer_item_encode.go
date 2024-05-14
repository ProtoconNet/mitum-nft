package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *TransferItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, rc string,
	nid uint64,
	cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	switch a, err := mitumbase.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	receiver, err := mitumbase.DecodeAddress(rc, enc)
	if err != nil {
		return err
	}
	it.receiver = receiver
	it.nft = nid
	it.currency = types.CurrencyID(cid)

	return nil
}
