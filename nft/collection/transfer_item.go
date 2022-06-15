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
	TransferItemType   = hint.Type("mitum-nft-transfer-item")
	TransferItemHint   = hint.NewHint(TransferItemType, "v0.0.1")
	TransferItemHinter = TransferItem{BaseHinter: hint.NewBaseHinter(TransferItemHint)}
)

type TransferItem struct {
	hint.BaseHinter
	receiver base.Address
	nft      nft.NFTID
	cid      currency.CurrencyID
}

func NewTransferItem(receiver base.Address, n nft.NFTID, cid currency.CurrencyID) TransferItem {
	return TransferItem{
		BaseHinter: hint.NewBaseHinter(TransferItemHint),
		receiver:   receiver,
		nft:        n,
		cid:        cid,
	}
}

func (it TransferItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.receiver.Bytes(),
		it.nft.Bytes(),
		it.cid.Bytes(),
	)
}

func (it TransferItem) IsValid([]byte) error {
	return isvalid.Check(nil, false, it.BaseHinter, it.receiver, it.nft, it.cid)
}

func (it TransferItem) Receiver() base.Address {
	return it.receiver
}

func (it TransferItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it TransferItem) NFT() nft.NFTID {
	return it.nft
}

func (it TransferItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it TransferItem) Rebuild() TransferItem {
	return it
}
