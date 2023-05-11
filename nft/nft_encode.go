package nft

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (n *NFT) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	bid []byte,
	ac bool,
	ow string,
	hs string,
	uri string,
	ap string,
	bcrs []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFT")

	n.BaseHinter = hint.NewBaseHinter(ht)
	n.active = ac
	n.hash = NFTHash(hs)
	n.uri = URI(uri)

	owner, err := base.DecodeAddress(ow, enc)
	if err != nil {
		return e(err, "")
	}
	n.owner = owner

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return e(err, "")
	}
	n.approved = approved

	if hinter, err := enc.Decode(bid); err != nil {
		return e(err, "")
	} else if id, ok := hinter.(NFTID); !ok {
		return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
	} else {
		n.id = id
	}

	if hinter, err := enc.Decode(bcrs); err != nil {
		return e(err, "")
	} else if sns, ok := hinter.(Signers); !ok {
		return e(util.ErrWrongType.Errorf("expected Signers, not %T", hinter), "")
	} else {
		n.creators = sns
	}

	return nil
}
