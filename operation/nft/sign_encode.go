package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *SignFact) unpack(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	sender, err := mitumbase.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]SignItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(SignItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected SignItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
