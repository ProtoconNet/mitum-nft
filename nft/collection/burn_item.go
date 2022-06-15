package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	BurnItemType   = hint.Type("mitum-nft-burn-item")
	BurnItemHint   = hint.NewHint(BurnItemType, "v0.0.1")
	BurnItemHinter = BurnItem{BaseHinter: hint.NewBaseHinter(BurnItemHint)}
)

type BurnItem struct {
	hint.BaseHinter
	nft nft.NFTID
	cid currency.CurrencyID
}

func NewBurnItem(n nft.NFTID, cid currency.CurrencyID) BurnItem {
	return BurnItem{
		BaseHinter: hint.NewBaseHinter(BurnItemHint),
		nft:        n,
		cid:        cid,
	}
}

func (it BurnItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.nft.Bytes(),
		it.cid.Bytes(),
	)
}

func (it BurnItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.nft, it.cid); err != nil {
		return err
	}

	return nil
}

func (it BurnItem) NFT() nft.NFTID {
	return it.nft
}

func (it BurnItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BurnItem) Rebuild() BurnItem {
	return it
}
