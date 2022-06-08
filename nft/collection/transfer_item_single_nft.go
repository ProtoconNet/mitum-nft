package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	TransferItemSingleNFTType   = hint.Type("mitum-nft-tranfer-single-nft")
	TransferItemSingleNFTHint   = hint.NewHint(TransferItemSingleNFTType, "v0.0.1")
	TransferItemSingleNFTHinter = TransferItemSingleNFT{
		BaseTransferItem: BaseTransferItem{
			BaseHinter: hint.NewBaseHinter(TransferItemSingleNFTHint),
		},
	}
)

type TransferItemSingleNFT struct {
	BaseTransferItem
}

func NewTransferItemSingleNFT(receiver base.Address, nftid nft.NFTID, cid currency.CurrencyID) TransferItemSingleNFT {
	return TransferItemSingleNFT{
		BaseTransferItem: NewBaseTransferItem(TransferItemSingleNFTHint, receiver, []nft.NFTID{nftid}, cid),
	}
}

func (it TransferItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseTransferItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it TransferItemSingleNFT) Rebuild() TransferItem {
	it.BaseTransferItem = it.BaseTransferItem.Rebuild().(BaseTransferItem)

	return it
}
