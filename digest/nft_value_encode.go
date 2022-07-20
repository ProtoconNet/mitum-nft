package digest

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (n *NFTValue) unpack(enc encoder.Encoder, bdm []byte, height base.Height) error {
	if bdm != nil {
		i, err := DecodeNFT(bdm, enc)
		if err != nil {
			return err
		}
		n.nft = i
	}

	n.height = height

	return nil
}
