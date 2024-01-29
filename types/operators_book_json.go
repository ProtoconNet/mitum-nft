package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type OperatorsBookJSONMarshaler struct {
	hint.BaseHinter
	Operators []base.Address `json:"operators"`
}

func (ob OperatorsBook) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperatorsBookJSONMarshaler{
		BaseHinter: ob.BaseHinter,
		Operators:  ob.operators,
	})
}

type OperatorsBookJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Operators []string  `json:"operators"`
}

func (ob *OperatorsBook) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of operators book")

	var u OperatorsBookJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return ob.unmarshal(enc, u.Hint, u.Operators)
}
