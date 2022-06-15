package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *TransferItem) unpack(
	enc encoder.Encoder,
	br base.AddressDecoder,
	bn []byte,
	cid string,
) error {
	receiver, err := br.Encode(enc)
	if err != nil {
		return err
	}
	it.receiver = receiver

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
