package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type SignItemJSONMarshaler struct {
	hint.BaseHinter
	Contract mitumbase.Address `json:"contract"`
	NFT      uint64            `json:"nft"`
	Currency types.CurrencyID  `json:"currency"`
}

func (it SignItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SignItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		NFT:        it.nft,
		Currency:   it.currency,
	})
}

type SignItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	NFT      uint64    `json:"nft"`
	Currency string    `json:"currency"`
}

func (it *SignItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u SignItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.NFT, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
