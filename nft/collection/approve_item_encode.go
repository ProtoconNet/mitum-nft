package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *ApproveItem) unpack(
	enc encoder.Encoder,
	bap base.AddressDecoder,
	bn []byte,
	cid string,
) error {
	approved, err := bap.Encode(enc)
	if err != nil {
		return err
	}
	it.approved = approved

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		it.nft = n
	}

	it.cid = currency.CurrencyID(cid)

	return nil
}
