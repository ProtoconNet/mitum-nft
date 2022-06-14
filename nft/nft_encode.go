package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (n *NFT) unpack(
	enc encoder.Encoder,
	bid []byte,
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

	hcrs, err := enc.DecodeSlice(bcrs)
	if err != nil {
		return err
	}
	crs := make([]RightHolder, len(hcrs))
	for i := range hcrs {
		r, ok := hcrs[i].(RightHolder)
		if !ok {
			return util.WrongTypeError.Errorf("not RightHolder; %T", hcrs[i])
		}
		crs[i] = r
	}
	n.creators = crs

	hcps, err := enc.DecodeSlice(bcps)
	if err != nil {
		return err
	}
	cps := make([]RightHolder, len(hcps))
	for i := range hcps {
		r, ok := hcps[i].(RightHolder)
		if !ok {
			return util.WrongTypeError.Errorf("not RightHolder; %T", hcps[i])
		}
		cps[i] = r
	}
	n.copyrighters = cps

	return nil
}
