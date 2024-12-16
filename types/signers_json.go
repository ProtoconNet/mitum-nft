package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type SignersJSONMarshaler struct {
	hint.BaseHinter
	Signers []Signer `json:"signers"`
}

func (sgns Signers) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SignersJSONMarshaler{
		BaseHinter: sgns.BaseHinter,
		Signers:    sgns.signers,
	})
}

type SignersJSONUnmarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	Signers json.RawMessage `json:"signers"`
}

func (sgns *Signers) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Signers")

	var u SignersJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return sgns.unpack(enc, u.Hint, u.Signers)
}
