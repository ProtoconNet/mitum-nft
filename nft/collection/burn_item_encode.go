package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BaseBurnItem) unpack(
	enc encoder.Encoder,
	collection string,
	bns []byte,
	cid string,
) error {
	it.collection = extensioncurrency.ContractID(collection)

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
	it.cid = currency.CurrencyID(cid)

	return nil
}
