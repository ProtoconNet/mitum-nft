package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MintItemSingleNFTType   = hint.Type("mitum-nft-mint-single-nft")
	MintItemSingleNFTHint   = hint.NewHint(MintItemSingleNFTType, "v0.0.1")
	MintItemSingleNFTHinter = MintItemSingleNFT{
		BaseMintItem: BaseMintItem{
			BaseHinter: hint.NewBaseHinter(MintItemSingleNFTHint),
		},
	}
)

type MintItemSingleNFT struct {
	BaseMintItem
}

func NewMintItemSingleNFT(collection extensioncurrency.ContractID, form MintForm, cid currency.CurrencyID) MintItemSingleNFT {
	return MintItemSingleNFT{
		BaseMintItem: NewBaseMintItem(MintItemSingleNFTHint, collection, []MintForm{form}, cid),
	}
}

func (it MintItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseMintItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.forms); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it MintItemSingleNFT) Rebuild() MintItem {
	it.BaseMintItem = it.BaseMintItem.Rebuild().(BaseMintItem)

	return it
}
