package nft

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *ApproveFact) unmarshal(enc encoder.Encoder, sd string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal ApproveFact")

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e(err, "")
	}

	items := make([]ApproveItem, len(hit))
	for i, hinter := range hit {
		item, ok := hinter.(ApproveItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected ApproveItem, not %T", hinter), "")
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
