package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type SignItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   mitumbase.Address `json:"contract"`
	Collection types.ContractID  `json:"collection"`
	NFT        uint64            `json:"nft"`
	Currency   types.CurrencyID  `json:"currency"`
}

func (it SignItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SignItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Collection: it.collection,
		NFT:        it.nft,
		Currency:   it.currency,
	})
}

type SignItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Collection string    `json:"collection"`
	NFT        uint64    `json:"nft"`
	Currency   string    `json:"currency"`
}

func (it *SignItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of SignItem")

	var u SignItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.NFT, u.Currency)
}
