package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (ob *OperatorsBook) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	bags []string,
) error {
	e := util.StringErrorFunc("failed to unmarshal operators book")

	ob.BaseHinter = hint.NewBaseHinter(ht)
	ob.collection = currencybase.ContractID(col)

	operators := make([]base.Address, len(bags))
	for i, bag := range bags {
		operator, err := base.DecodeAddress(bag, enc)
		if err != nil {
			return e(err, "")
		}
		operators[i] = operator
	}
	ob.operators = operators

	return nil
}
