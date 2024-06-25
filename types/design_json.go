package types

import (
	"encoding/json"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Contract mitumbase.Address `json:"contract"`
	Creator  mitumbase.Address `json:"creator"`
	Active   bool              `json:"active"`
	Policy   BasePolicy        `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Contract:   de.contract,
		Creator:    de.creator,
		Active:     de.active,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Contract string          `json:"contract"`
	Creator  string          `json:"creator"`
	Active   bool            `json:"active"`
	Policy   json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Design")

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return de.unpack(enc, u.Hint, u.Contract, u.Creator, u.Active, u.Policy)
}
