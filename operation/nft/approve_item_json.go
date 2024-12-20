package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address     `json:"contract"`
	Approved base.Address     `json:"approved"`
	NFTIdx   uint64           `json:"nft_idx"`
	Currency types.CurrencyID `json:"currency"`
}

func (it ApproveItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Approved:   it.approved,
		NFTIdx:     it.nftIdx,
		Currency:   it.currency,
	})
}

type ApproveItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Approved string    `json:"approved"`
	NFTIdx   uint64    `json:"nft_idx"`
	Currency string    `json:"currency"`
}

func (it *ApproveItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u ApproveItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Approved, u.NFTIdx, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
