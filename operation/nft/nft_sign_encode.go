package nft

import (
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *NFTSignFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTSignFact")

	sender, err := mitumbase.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return err
	}

	items := make([]NFTSignItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(NFTSignItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected SignItem, not %T", hinter), "")
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
