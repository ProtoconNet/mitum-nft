package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsApproveItemMultiNFTs = 10

var (
	ApproveItemMultiNFTsType   = hint.Type("mitum-nft-approve-multi-nfts")
	ApproveItemMultiNFTsHint   = hint.NewHint(ApproveItemMultiNFTsType, "v0.0.1")
	ApproveItemMultiNFTsHinter = ApproveItemMultiNFTs{
		BaseApproveItem: BaseApproveItem{
			BaseHinter: hint.NewBaseHinter(ApproveItemMultiNFTsHint),
		},
	}
)

type ApproveItemMultiNFTs struct {
	BaseApproveItem
}

func NewApproveItemMultiNFTs(approved base.Address, nfts []nft.NFTID, cid currency.CurrencyID) ApproveItemMultiNFTs {
	return ApproveItemMultiNFTs{
		BaseApproveItem: NewBaseApproveItem(ApproveItemMultiNFTsHint, approved, nfts, cid),
	}
}

func (it ApproveItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseApproveItem.IsValid(nil); err != nil {
		return err
	}

	if l := len(it.nfts); l > MaxNFTsApproveItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", l, MaxNFTsApproveItemMultiNFTs)
	}

	return nil
}

func (it ApproveItemMultiNFTs) Rebuild() ApproveItem {
	it.BaseApproveItem = it.BaseApproveItem.Rebuild().(BaseApproveItem)

	return it
}
