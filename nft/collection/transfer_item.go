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
	receiver base.Address
	nfts     []nft.NFTID
	cid      currency.CurrencyID
}

func NewBaseTransferItem(ht hint.Hint, receiver base.Address, nfts []nft.NFTID, cid currency.CurrencyID) BaseTransferItem {
	return BaseTransferItem{
		BaseHinter: hint.NewBaseHinter(ht),
		receiver:   receiver,
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseTransferItem) Bytes() []byte {
	bns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		bns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.receiver.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(bns...),
	)
}

func (it BaseTransferItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.receiver, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseTransferItem")
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}
		n := it.nfts[i]
		if _, found := foundNFT[n]; found {
			return errors.Errorf("duplicated nft found; %s", n)
		}
		foundNFT[n] = true
	}

	return nil
}

func (it BaseTransferItem) Receiver() base.Address {
	return it.receiver
}

func (it BaseTransferItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
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
