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

type BaseTransferItem struct {
	hint.BaseHinter
	from base.Address
	to   base.Address
	nfts []nft.NFTID
	cid  currency.CurrencyID
}

func NewBaseTransferItem(ht hint.Hint, from base.Address, to base.Address, nfts []nft.NFTID, cid currency.CurrencyID) BaseTransferItem {
	return BaseTransferItem{
		BaseHinter: hint.NewBaseHinter(ht),
		from:       from,
		to:         to,
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseTransferItem) Bytes() []byte {
	ns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		ns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.from.Bytes(),
		it.to.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (it BaseTransferItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.from, it.to, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseTransferItem")
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

func (it BaseTransferItem) From() base.Address {
	return it.from
}

func (it BaseTransferItem) To() base.Address {
	return it.to
}

func (it BaseTransferItem) Addresses() []base.Address {
	as := make([]base.Address, 2)
	as[0] = it.From()
	as[1] = it.To()
	return as
}

func (it BaseTransferItem) NFTs() []nft.NFTID {
	return it.nfts
}

func (it BaseTransferItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseTransferItem) Rebuild() TransferItem {
	nfts := make([]nft.NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
