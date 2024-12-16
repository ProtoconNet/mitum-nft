package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (de *Design) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	pAdr string,
	crAdr string,
	active bool,
	bPcy []byte,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)
	de.active = active

	contract, err := base.DecodeAddress(pAdr, enc)
	if err != nil {
		return err
	}
	de.contract = contract

	creator, err := base.DecodeAddress(crAdr, enc)
	if err != nil {
		return err
	}
	de.creator = creator

	if hinter, err := enc.Decode(bPcy); err != nil {
		return err
	} else if po, ok := hinter.(BasePolicy); !ok {
		return errors.Errorf("expected BasePolicy, not %T", hinter)
	} else {
		de.policy = po
	}

	return nil
}
