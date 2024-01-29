package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TransferItemJSONMarshaler struct {
	hint.BaseHinter
	Contract mitumbase.Address `json:"contract"`
	Receiver mitumbase.Address `json:"receiver"`
	NFTidx   uint64            `json:"nft"`
	Currency types.CurrencyID  `json:"currency"`
}

func (it TransferItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		NFTidx:     it.nft,
		Currency:   it.currency,
	})
}

type TransferItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Receiver string    `json:"receiver"`
	NFTidx   uint64    `json:"nft"`
	Currency string    `json:"currency"`
}

func (it *TransferItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TransferItem")

	var u TransferItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Receiver, u.NFTidx, u.Currency)
}
