package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BaseApproveItem) unpack(
	enc encoder.Encoder,
	bApproved base.AddressDecoder,
	bNFTs []byte,
	cid string,
) error {
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

	it.approved = approved
	it.nfts = nfts
	it.cid = currency.CurrencyID(cid)

	return nil
}
