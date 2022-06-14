package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseApproveItem struct {
	hint.BaseHinter
	approved base.Address
	nfts     []nft.NFTID
	cid      currency.CurrencyID
}

func NewBaseApproveItem(ht hint.Hint, approved base.Address, nfts []nft.NFTID, cid currency.CurrencyID) BaseApproveItem {
	return BaseApproveItem{
		BaseHinter: hint.NewBaseHinter(ht),
		approved:   approved,
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseApproveItem) Bytes() []byte {
	bns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		bns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.approved.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(bns...),
	)
}

func (it BaseApproveItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.approved, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseApproveItem")
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

func (it BaseApproveItem) Approved() base.Address {
	return it.approved
}

func (it BaseApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it BaseApproveItem) NFTs() []nft.NFTID {
	return it.nfts
}

func (it BaseApproveItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseApproveItem) Rebuild() ApproveItem {
	nfts := make([]nft.NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
