package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *ApproveFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bApproved base.AddressDecoder,
	bNFTs []byte,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	approved, err := bApproved.Encode(enc)
	if err != nil {
		return err
	}

	hNFTs, err := enc.DecodeSlice(bNFTs)
	if err != nil {
		return err
	}

	nfts := make([]nft.NFTID, len(hNFTs))
	for i := range hNFTs {
		j, ok := hNFTs[i].(nft.NFTID)
		if !ok {
			return util.WrongTypeError.Errorf("not NFTID; %T", hNFTs[i])
		}

		nfts[i] = j
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.approved = approved
	fact.nfts = nfts
	fact.cid = currency.CurrencyID(cid)

	return nil
}
