package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (n *NFT) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	id uint64,
	ac bool,
	ow string,
	hs string,
	uri string,
	ap string,
	bcrs []byte,
) error {
	n.BaseHinter = hint.NewBaseHinter(ht)
	n.active = ac
	n.hash = NFTHash(hs)
	n.uri = URI(uri)

	owner, err := base.DecodeAddress(ow, enc)
	if err != nil {
		return err
	}
	n.owner = owner

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return err
	}
	n.approved = approved
	n.id = id

	if hinter, err := enc.Decode(bcrs); err != nil {
		return err
	} else if sns, ok := hinter.(Signers); !ok {
		return errors.Errorf("expected Signers, not %T", hinter)
	} else {
		n.creators = sns
	}

	return nil
}
