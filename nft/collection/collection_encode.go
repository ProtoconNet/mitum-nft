package collection

import (
	"net/url"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (p *Policy) unpack(
	enc encoder.Encoder,
	name string,
	bCreator base.AddressDecoder,
	royalty uint,
	_uri string,
) error {
	p.name = CollectionName(name)

	creator, err := bCreator.Encode(enc)
	if err != nil {
		return err
	}
	p.creator = creator

	p.royalty = nft.PaymentParameter(royalty)

	if uri, err := url.Parse(_uri); err != nil {
		return err
	} else {
		p.uri = *uri
	}

	return nil
}
