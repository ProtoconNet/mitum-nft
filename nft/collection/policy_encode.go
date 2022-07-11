package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (p *CollectionPolicy) unpack(
	enc encoder.Encoder,
	name string,
	royalty uint,
	uri string,
	bws []base.AddressDecoder,
) error {
	p.name = CollectionName(name)
	p.royalty = nft.PaymentParameter(royalty)
	p.uri = nft.URI(uri)

	whites := make([]base.Address, len(bws))
	for i := range bws {
		if white, err := bws[i].Encode(enc); err != nil {
			return err
		} else {
			whites[i] = white
		}
	}
	p.whites = whites

	return nil
}
