package nft

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *TransferFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	e := util.StringError("failed to unmarshal TransferFact")

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]TransferItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(TransferItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected TransferItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
