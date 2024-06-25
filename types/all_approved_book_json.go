package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type OperatorsBookJSONMarshaler struct {
	hint.BaseHinter
	AllApproved []base.Address `json:"all_approved"`
}

func (ob AllApprovedBook) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperatorsBookJSONMarshaler{
		BaseHinter:  ob.BaseHinter,
		AllApproved: ob.allApproved,
	})
}

type OperatorsBookJSONUnmarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	AllApproved []string  `json:"all_approved"`
}

func (ob *AllApprovedBook) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("decode json of all-approved book")

	var u OperatorsBookJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return ob.unpack(enc, u.Hint, u.AllApproved)
}
