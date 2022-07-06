package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (n *NFT) unpack(
	enc encoder.Encoder,
	bid []byte,
	active bool,
	bo base.AddressDecoder,
	hash string,
	uri string,
	bap base.AddressDecoder,
	bcrs []byte,
	bcps []byte,
) error {
	if hinter, err := enc.Decode(bid); err != nil {
		return err
	} else if id, ok := hinter.(NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		n.id = id
	}

	n.active = active

	owner, err := bo.Encode(enc)
	if err != nil {
		return err
	}
	n.owner = owner

	approved, err := bap.Encode(enc)
	if err != nil {
		return err
	}
	n.approved = approved

	n.uri = URI(uri)
	n.hash = NFTHash(hash)

	if hinter, err := enc.Decode(bcrs); err != nil {
		return err
	} else if sns, ok := hinter.(Signers); !ok {
		return util.WrongTypeError.Errorf("not Signers; %T", hinter)
	} else {
		n.creators = sns
	}

	if hinter, err := enc.Decode(bcps); err != nil {
		return err
	} else if sns, ok := hinter.(Signers); !ok {
		return util.WrongTypeError.Errorf("not Signer; %T", hinter)
	} else {
		n.copyrighters = sns
	}

	return nil
}
