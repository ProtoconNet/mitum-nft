package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	ApproveItemType   = hint.Type("mitum-nft-approve-item")
	ApproveItemHint   = hint.NewHint(ApproveItemType, "v0.0.1")
	ApproveItemHinter = ApproveItem{BaseHinter: hint.NewBaseHinter(ApproveItemHint)}
)

type ApproveItem struct {
	hint.BaseHinter
	approved base.Address
	nft      nft.NFTID
	cid      currency.CurrencyID
}

func NewApproveItem(approved base.Address, n nft.NFTID, cid currency.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		approved:   approved,
		nft:        n,
		cid:        cid,
	}
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.approved.Bytes(),
		it.nft.Bytes(),
		it.cid.Bytes(),
	)
}

func (it ApproveItem) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		it.BaseHinter,
		it.approved,
		it.nft,
		it.cid); err != nil {
		return err
	}
	return nil
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFT() nft.NFTID {
	return it.nft
}

func (it ApproveItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it ApproveItem) Rebuild() ApproveItem {
	return it
}
