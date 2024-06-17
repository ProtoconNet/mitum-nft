package types

import (
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (de *Design) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	pAdr string,
	crAdr string,
	active bool,
	bPcy []byte,
) error {
	e := util.StringError("failed to unmarshal Design")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.active = active

	contract, err := mitumbase.DecodeAddress(pAdr, enc)
	if err != nil {
		return e.Wrap(err)
	}
	de.contract = contract

	creator, err := mitumbase.DecodeAddress(crAdr, enc)
	if err != nil {
		return e.Wrap(err)
	}
	de.creator = creator

	if hinter, err := enc.Decode(bPcy); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(BasePolicy); !ok {
		return e.Wrap(errors.Errorf("expected BasePolicy, not %T", hinter))
	} else {
		de.policy = po
	}

	return nil
}
