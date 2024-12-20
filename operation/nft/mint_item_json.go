package nft

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type MintItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address      `json:"contract"`
	Receiver base.Address      `json:"receiver"`
	Hash     types.NFTHash     `json:"hash"`
	Uri      types.URI         `json:"uri"`
	Creators types.Signers     `json:"creators"`
	Currency ctypes.CurrencyID `json:"currency"`
}

func (it MintItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		Hash:       it.hash,
		Uri:        it.uri,
		Creators:   it.creators,
		Currency:   it.currency,
	})
}

type MintItemJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Contract string          `json:"contract"`
	Receiver string          `json:"receiver"`
	Hash     string          `json:"hash"`
	Uri      string          `json:"uri"`
	Creators json.RawMessage `json:"creators"`
	Currency string          `json:"currency"`
}

func (it *MintItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u MintItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Receiver, u.Hash, u.Uri, u.Creators, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
