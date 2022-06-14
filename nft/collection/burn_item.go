package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseBurnItem struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	nfts       []nft.NFTID
	cid        currency.CurrencyID
}

func NewBaseBurnItem(ht hint.Hint, collection extensioncurrency.ContractID, nfts []nft.NFTID, cid currency.CurrencyID) BaseBurnItem {
	return BaseBurnItem{
		BaseHinter: hint.NewBaseHinter(ht),
		collection: collection,
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseBurnItem) Bytes() []byte {
	bns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		bns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.collection.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(bns...),
	)
}

func (it BaseBurnItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.collection, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseBurnItem")
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}

		n := it.nfts[i]
		if _, found := foundNFT[n]; found {
			return errors.Errorf("duplicated nft found; %q", n)
		}

		foundNFT[n] = true
	}

	return nil
}

func (it BaseBurnItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it BaseBurnItem) NFTs() []nft.NFTID {
	return it.nfts
}

func (it BaseBurnItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseBurnItem) Rebuild() BurnItem {
	nfts := make([]nft.NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
