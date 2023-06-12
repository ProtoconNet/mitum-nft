package nft

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *MintFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal MintFact")

	switch sender, err := base.DecodeAddress(sd, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = sender
	}

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return e(err, "")
	}

	items := make([]MintItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(MintItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected MintItem, not %T", hinter), "")
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
