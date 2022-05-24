package nft

import (
	"net/url"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (nid *NFTID) unpack(
	enc encoder.Encoder,
	collection string,
	idx currency.Big,
) error {
	nid.collection = extension.ContractID(collection)
	nid.idx = idx

	return nil
}

func (nft *NFT) unpack(
	enc encoder.Encoder,
	bId []byte,
	bOwner base.AddressDecoder,
	hash string,
	_uri string,
	bApproved base.AddressDecoder,
	bCopyrighter base.AddressDecoder,
) error {
	if hinter, err := enc.Decode(bId); err != nil {
		return err
	} else if id, ok := hinter.(NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		nft.id = id
	}

	owner, err := bOwner.Encode(enc)
	if err != nil {
		return err
	}
	nft.owner = owner

	approved, err := bApproved.Encode(enc)
	if err != nil {
		return err
	}
	nft.approved = approved

	if uri, err := url.Parse(_uri); err != nil {
		return err
	} else {
		nft.uri = *uri
	}
	nft.hash = NFTHash(hash)

	copyrighter, err := bCopyrighter.Encode(enc)
	if err != nil {
		return err
	}
	nft.copyrighter = copyrighter

	return nil
}
