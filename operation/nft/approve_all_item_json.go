package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveAllItemJSONMarshaler struct {
	hint.BaseHinter
	Contract mitumbase.Address `json:"contract"`
	Approved mitumbase.Address `json:"approved"`
	Mode     ApproveAllMode    `json:"mode"`
	Currency types.CurrencyID  `json:"currency"`
}

func (it ApproveAllItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveAllItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Approved:   it.approved,
		Mode:       it.mode,
		Currency:   it.currency,
	})
}

type ApproveAllItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Approved string    `json:"approved"`
	Mode     string    `json:"mode"`
	Currency string    `json:"currency"`
}

func (it *ApproveAllItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u ApproveAllItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unmarshal(enc, u.Hint, u.Contract, u.Approved, u.Mode, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
