package timestamp

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Service  extensioncurrency.ContractID `json:"service"`
	Projects []string                     `json:"projects"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Service:    de.service,
		Projects:   de.projects,
	})
}

type DesignJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Service  string    `json:"service"`
	Projects []string  `json:"projects"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Design")

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return de.unmarshal(enc, u.Hint, u.Service, u.Projects)
}
