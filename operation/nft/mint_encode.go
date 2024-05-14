package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *MintFact) unpack(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	switch sender, err := base.DecodeAddress(sd, enc); {
	case err != nil:
		return err
	default:
		fact.sender = sender
	}

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]MintItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(MintItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected MintItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
