package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
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
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	receiver, err := base.DecodeAddress(rc, enc)
	if err != nil {
		return err
	}
	it.receiver = receiver
	it.nftIdx = nid
	it.currency = types.CurrencyID(cid)

	return nil
}
