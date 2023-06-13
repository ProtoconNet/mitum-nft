package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	pr string,
	cr string,
	sb string,
	ac bool,
	bpo []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal Design")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.collection = types.ContractID(sb)
	de.active = ac

	parent, err := mitumbase.DecodeAddress(pr, enc)
	if err != nil {
		return e(err, "")
	}
	de.parent = parent

	creator, err := mitumbase.DecodeAddress(cr, enc)
	if err != nil {
		return e(err, "")
	}
	de.creator = creator

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(BasePolicy); !ok {
		return e(util.ErrWrongType.Errorf("expected BasePolicy, not %T", hinter), "")
	} else {
		de.policy = po
	}

	return nil
}
