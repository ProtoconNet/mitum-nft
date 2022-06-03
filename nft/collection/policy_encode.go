package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/util/encoder"
)

func (p *CollectionPolicy) unpack(
	enc encoder.Encoder,
	name string,
	royalty uint,
	uri string,
) error {
	p.name = CollectionName(name)
	p.royalty = nft.PaymentParameter(royalty)
	p.uri = nft.URI(uri)

	return nil
}
