package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsMintItemMultiNFTs = 10

var (
	MintItemMultiNFTsType   = hint.Type("mitum-nft-mint-multi-nfts")
	MintItemMultiNFTsHint   = hint.NewHint(MintItemMultiNFTsType, "v0.0.1")
	MintItemMultiNFTsHinter = MintItemMultiNFTs{
		BaseMintItem: BaseMintItem{
			BaseHinter: hint.NewBaseHinter(MintItemMultiNFTsHint),
		},
	}
)

type MintItemMultiNFTs struct {
	BaseMintItem
}

func NewMintItemMultiNFTs(collection extensioncurrency.ContractID, forms []MintForm, cid currency.CurrencyID) MintItemMultiNFTs {
	return MintItemMultiNFTs{
		BaseMintItem: NewBaseMintItem(MintItemMultiNFTsHint, collection, forms, cid),
	}
}

func (it MintItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseMintItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.forms); n > MaxNFTsMintItemMultiNFTs {
		return isvalid.InvalidError.Errorf("forms over allowed; %d > %d", n, MaxNFTsMintItemMultiNFTs)
	}

	return nil
}

func (it MintItemMultiNFTs) Rebuild() MintItem {
	it.BaseMintItem = it.BaseMintItem.Rebuild().(BaseMintItem)

	return it
}
