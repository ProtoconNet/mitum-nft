package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BaseTransferItem) unpack(
	enc encoder.Encoder,
	br base.AddressDecoder,
	bns []byte,
	cid string,
) error {
	receiver, err := br.Encode(enc)
	if err != nil {
		return err
	}

	hns, err := enc.DecodeSlice(bns)
	if err != nil {
		return err
	}

	nfts := make([]nft.NFTID, len(hns))
	for i := range hns {
		n, ok := hns[i].(nft.NFTID)
		if !ok {
			return util.WrongTypeError.Errorf("not NFTID; %T", hns[i])
		}

		nfts[i] = n
	}
	it.nfts = nfts

	it.receiver = receiver
	it.cid = currency.CurrencyID(cid)

	return nil
}
