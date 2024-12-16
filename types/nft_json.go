package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTJSONMarshaler struct {
	hint.BaseHinter
	ID       uint64       `json:"nft_idx"`
	Active   bool         `json:"active"`
	Owner    base.Address `json:"owner"`
	Hash     NFTHash      `json:"hash"`
	URI      URI          `json:"uri"`
	Approved base.Address `json:"approved"`
	Creators Signers      `json:"creators"`
}

func (n NFT) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTJSONMarshaler{
		BaseHinter: n.BaseHinter,
		ID:         n.id,
		Active:     n.active,
		Owner:      n.owner,
		Hash:       n.hash,
		URI:        n.uri,
		Approved:   n.approved,
		Creators:   n.creators,
	})
}

type NFTJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	ID       uint64          `json:"nft_idx"`
	Active   bool            `json:"active"`
	Owner    string          `json:"owner"`
	Hash     string          `json:"hash"`
	URI      string          `json:"uri"`
	Approved string          `json:"approved"`
	Creators json.RawMessage `json:"creators"`
}

func (n *NFT) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of NFT")

	var u NFTJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return n.unpack(enc, u.Hint, u.ID, u.Active, u.Owner, u.Hash, u.URI, u.Approved, u.Creators)
}
