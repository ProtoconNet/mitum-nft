package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
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
	NFTIdx   uint64            `json:"nft_idx"`
	Currency types.CurrencyID  `json:"currency"`
}

func (it TransferItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		NFTIdx:     it.nftIdx,
		Currency:   it.currency,
	})
}

type TransferItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Receiver string    `json:"receiver"`
	NFTIdx   uint64    `json:"nft_idx"`
	Currency string    `json:"currency"`
}

func (it *TransferItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u TransferItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Receiver, u.NFTIdx, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
