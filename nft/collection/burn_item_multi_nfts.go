package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsBurnItemMultiNFTs = 10

var (
	BurnItemMultiNFTsType   = hint.Type("mitum-nft-burn-multi-nfts")
	BurnItemMultiNFTsHint   = hint.NewHint(BurnItemMultiNFTsType, "v0.0.1")
	BurnItemMultiNFTsHinter = BurnItemMultiNFTs{
		BaseBurnItem: BaseBurnItem{
			BaseHinter: hint.NewBaseHinter(BurnItemMultiNFTsHint),
		},
	}
)

type BurnItemMultiNFTs struct {
	BaseBurnItem
}

func NewBurnItemMultiNFTs(collection extensioncurrency.ContractID, nfts []nft.NFTID, cid currency.CurrencyID) BurnItemMultiNFTs {
	return BurnItemMultiNFTs{
		BaseBurnItem: NewBaseBurnItem(BurnItemMultiNFTsHint, collection, nfts, cid),
	}
}

func (it BurnItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseBurnItem.IsValid(nil); err != nil {
		return err
	}

	if l := len(it.nfts); l > MaxNFTsBurnItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", l, MaxNFTsBurnItemMultiNFTs)
	}

	return nil
}

func (it BurnItemMultiNFTs) Rebuild() BurnItem {
	it.BaseBurnItem = it.BaseBurnItem.Rebuild().(BaseBurnItem)

	return it
}
