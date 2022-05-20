package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (cudp *CollectionPolicy) unpack(
	enc encoder.Encoder,
	symbol string,
	name string,
	bCreator base.AddressDecoder,
	royalty uint,
	uri string,
) error {
	cudp.symbol = nft.Symbol(symbol)
	cudp.name = CollectionName(name)

	creator, err := bCreator.Encode(enc)
	if err != nil {
		return err
	}
	cudp.creator = creator

	cudp.royalty = nft.PaymentParameter(royalty)
	cudp.uri = CollectionUri(uri)

	return nil
}
