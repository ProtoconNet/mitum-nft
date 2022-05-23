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
	ns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		ns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.approved.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (it BaseApproveItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.approved, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseApproveItem")
	}

	foundNFT := map[string]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}
		nft := it.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return errors.Errorf("duplicated nft found; %s", nft)
		}
		foundNFT[nft] = true
	}

	return nil
}

func (it BaseApproveItem) Approved() base.Address {
	return it.approved
}

func (it BaseApproveItem) Addresses() []base.Address {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as
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
