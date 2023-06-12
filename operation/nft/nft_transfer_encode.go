package nft

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *NFTTransferFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTTransferFact")

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]NFTTransferItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(NFTTransferItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected NFTTransferItem, not %T", hinter), "")
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
