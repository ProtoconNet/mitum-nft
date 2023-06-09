package nft

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Parent     base.Address            `json:"parent"`
	Creator    base.Address            `json:"creator"`
	Collection currencybase.ContractID `json:"collection"`
	Active     bool                    `json:"active"`
	Policy     BasePolicy              `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Parent:     de.parent,
		Creator:    de.creator,
		Collection: de.collection,
		Active:     de.active,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Parent     string          `json:"parent"`
	Creator    string          `json:"creator"`
	Collection string          `json:"collection"`
	Active     bool            `json:"active"`
	Policy     json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Design")

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return de.unmarshal(enc, u.Hint, u.Parent, u.Creator, u.Collection, u.Active, u.Policy)
}
