package collection

import (
	"net/url"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/encoder"
)

func (p *Policy) unpack(
	enc encoder.Encoder,
	name string,
	royalty uint,
	_uri string,
	_limit string,
) error {
	p.name = CollectionName(name)

	p.royalty = nft.PaymentParameter(royalty)

	if uri, err := url.Parse(_uri); err != nil {
		return err
	} else {
		p.uri = *uri
	}

	if limit, err := currency.NewBigFromString(_limit); err != nil {
		return err
	} else {
		p.limit = limit
	}

	return nil
}
