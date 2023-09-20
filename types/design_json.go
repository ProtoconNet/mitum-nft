package types

import (
	"encoding/json"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Parent  mitumbase.Address `json:"parent"`
	Creator mitumbase.Address `json:"creator"`
	Active  bool              `json:"active"`
	Policy  BasePolicy        `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Parent:     de.parent,
		Creator:    de.creator,
		Active:     de.active,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	Parent  string          `json:"parent"`
	Creator string          `json:"creator"`
	Active  bool            `json:"active"`
	Policy  json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Design")

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return de.unmarshal(enc, u.Hint, u.Parent, u.Creator, u.Active, u.Policy)
}
