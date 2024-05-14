package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var TransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type TransferItem struct {
	hint.BaseHinter
	contract mitumbase.Address
	receiver mitumbase.Address
	nft      uint64
	currency types.CurrencyID
}

func NewTransferItem(contract mitumbase.Address, receiver mitumbase.Address, nft uint64, currency types.CurrencyID) TransferItem {
	return TransferItem{
		BaseHinter: hint.NewBaseHinter(TransferItemHint),
		contract:   contract,
		receiver:   receiver,
		nft:        nft,
		currency:   currency,
	}
}

func (it TransferItem) IsValid([]byte) error {
	if it.receiver.Equal(it.contract) {
		return errors.Errorf("receiver is same with contract")
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.receiver,
		it.currency,
	)
}

func (it TransferItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		util.Uint64ToBytes(it.nft),
		it.currency.Bytes(),
	)
}

func (it TransferItem) Contract() mitumbase.Address {
	return it.contract
}

func (it TransferItem) Receiver() mitumbase.Address {
	return it.receiver
}

func (it TransferItem) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it TransferItem) NFT() uint64 {
	return it.nft
}

func (it TransferItem) Currency() types.CurrencyID {
	return it.currency
}
