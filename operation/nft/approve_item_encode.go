package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *ApproveItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	cAdr, appr string,
	idx uint64,
	cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = types.CurrencyID(cid)
	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	approved, err := base.DecodeAddress(appr, enc)
	if err != nil {
		return err
	}
	it.approved = approved
	it.nftIdx = idx

	return nil
}
