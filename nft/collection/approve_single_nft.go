package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	ApproveItemSingleNFTType   = hint.Type("mitum-nft-approve-single-nft")
	ApproveItemSingleNFTHint   = hint.NewHint(ApproveItemSingleNFTType, "v0.0.1")
	ApproveItemSingleNFTHinter = ApproveItemSingleNFT{
		BaseApproveItem: BaseApproveItem{
			BaseHinter: hint.NewBaseHinter(ApproveItemSingleNFTHint),
		},
	}
)

type ApproveItemSingleNFT struct {
	BaseApproveItem
}

func NewApproveItemSingleNFT(approved base.Address, n nft.NFTID, cid currency.CurrencyID) ApproveItemSingleNFT {
	return ApproveItemSingleNFT{
		BaseApproveItem: NewBaseApproveItem(ApproveItemSingleNFTHint, approved, []nft.NFTID{n}, cid),
	}
}

func (it ApproveItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseApproveItem.IsValid(nil); err != nil {
		return err
	}

	if l := len(it.nfts); l != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", l)
	}

	return nil
}

func (it ApproveItemSingleNFT) Rebuild() ApproveItem {
	it.BaseApproveItem = it.BaseApproveItem.Rebuild().(BaseApproveItem)

	return it
}
