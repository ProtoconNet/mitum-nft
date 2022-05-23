package collection

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *ApproveFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bItems []byte,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	hits, err := enc.DecodeSlice(bItems)
	if err != nil {
		return err
	}

	its := make([]ApproveItem, len(hits))
	for i := range hits {
		j, ok := hits[i].(ApproveItem)
		if !ok {
			return util.WrongTypeError.Errorf("not ApproveItem; %T", hits[i])
		}

		its[i] = j
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.items = its

	return nil
}
