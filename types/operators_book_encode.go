package types

import (
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (ob *OperatorsBook) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	oprs []string,
) error {
	e := util.StringError("failed to unmarshal operators book")

	ob.BaseHinter = hint.NewBaseHinter(ht)

	operators := make([]mitumbase.Address, len(oprs))
	for i, bag := range oprs {
		operator, err := mitumbase.DecodeAddress(bag, enc)
		if err != nil {
			return e.Wrap(err)
		}
		operators[i] = operator
	}
	ob.operators = operators

	return nil
}
