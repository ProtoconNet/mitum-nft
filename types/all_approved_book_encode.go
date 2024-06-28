package types

import (
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (ob *AllApprovedBook) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	oprs []string,
) error {
	ob.BaseHinter = hint.NewBaseHinter(ht)

	operators := make([]mitumbase.Address, len(oprs))
	for i, bag := range oprs {
		operator, err := mitumbase.DecodeAddress(bag, enc)
		if err != nil {
			return err
		}
		operators[i] = operator
	}
	ob.allApproved = operators

	return nil
}
