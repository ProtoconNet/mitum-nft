package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	BurnItemSingleNFTType   = hint.Type("mitum-nft-burn-single-nft")
	BurnItemSingleNFTHint   = hint.NewHint(BurnItemSingleNFTType, "v0.0.1")
	BurnItemSingleNFTHinter = BurnItemSingleNFT{
		BaseBurnItem: BaseBurnItem{
			BaseHinter: hint.NewBaseHinter(BurnItemSingleNFTHint),
		},
	}
)

type BurnItemSingleNFT struct {
	BaseBurnItem
}

func NewBurnItemSingleNFT(collection extensioncurrency.ContractID, n nft.NFTID, cid currency.CurrencyID) BurnItemSingleNFT {
	return BurnItemSingleNFT{
		BaseBurnItem: NewBaseBurnItem(BurnItemSingleNFTHint, collection, []nft.NFTID{n}, cid),
	}
}

func (it BurnItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseBurnItem.IsValid(nil); err != nil {
		return err
	}

	if l := len(it.nfts); l != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", l)
	}

	return nil
}

func (it BurnItemSingleNFT) Rebuild() BurnItem {
	it.BaseBurnItem = it.BaseBurnItem.Rebuild().(BaseBurnItem)

	return it
}
