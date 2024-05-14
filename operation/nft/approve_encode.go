package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *ApproveFact) unmarshal(enc encoder.Encoder, sd string, bit []byte) error {
	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return err
	}

	items := make([]ApproveItem, len(hit))
	for i, hinter := range hit {
		item, ok := hinter.(ApproveItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected ApproveItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
