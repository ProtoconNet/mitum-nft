package collection

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *DelegateFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bs base.AddressDecoder,
	bits []byte,
) error {
	sender, err := bs.Encode(enc)
	if err != nil {
		return err
	}

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]DelegateItem, len(hits))
	for i := range hits {
		item, ok := hits[i].(DelegateItem)
		if !ok {
			return util.WrongTypeError.Errorf("not DelegateItem; %T", hits[i])
		}

		items[i] = item
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.items = items

	return nil
}
