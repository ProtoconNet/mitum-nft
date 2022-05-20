package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsTransferItemMultiNFTs = 10

var (
	TransferItemMultiNFTsType   = hint.Type("mitum-nft-transfer-multi-nfts")
	TransferItemMultiNFTsHint   = hint.NewHint(TransferItemMultiNFTsType, "v0.0.1")
	TransferItemMultiNFTsHinter = TransferItemMultiNFTs{
		BaseTransferItem: BaseTransferItem{
			BaseHinter: hint.NewBaseHinter(TransferItemMultiNFTsHint),
		},
	}
)

type TransferItemMultiNFTs struct {
	BaseTransferItem
}

func NewTransferItemMultiNFTs(from base.Address, to base.Address, nfts []nft.NFTID, cid currency.CurrencyID) TransferItemMultiNFTs {
	return TransferItemMultiNFTs{
		BaseTransferItem: NewBaseTransferItem(TransferItemMultiNFTsHint, from, to, nfts, cid),
	}
}

func (it TransferItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseTransferItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n > MaxNFTsTransferItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", n, MaxNFTsTransferItemMultiNFTs)
	}

	return nil
}

func (it TransferItemMultiNFTs) Rebuild() TransferItem {
	it.BaseTransferItem = it.BaseTransferItem.Rebuild().(BaseTransferItem)

	return it
}
