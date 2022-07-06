package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (signer *Signer) unpack(
	enc encoder.Encoder,
	ba base.AddressDecoder,
	share uint,
	signed bool,
) error {
	a, err := ba.Encode(enc)
	if err != nil {
		return err
	}
	signer.account = a

	signer.share = share
	signer.signed = signed

	return nil
}
