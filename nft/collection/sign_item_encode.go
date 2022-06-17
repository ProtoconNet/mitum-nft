package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *SignItem) unpack(
	enc encoder.Encoder,
	q string,
	bn []byte,
	cid string,
) error {

	it.qualification = Qualification(q)

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
