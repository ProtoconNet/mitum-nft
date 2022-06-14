package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (r *RightHolder) unpack(
	enc encoder.Encoder,
	ba base.AddressDecoder,
	signed bool,
	clue string,
) error {
	a, err := ba.Encode(enc)
	if err != nil {
		return err
	}
	r.account = a

	r.signed = signed
	r.clue = clue

	return nil
}
