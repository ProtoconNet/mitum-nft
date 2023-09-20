package types

import (
	"bytes"
	"sort"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxOperators = 10

var OperatorsBookHint = hint.MustNewHint("mitum-nft-operator-book-v0.0.1")

type OperatorsBook struct {
	hint.BaseHinter
	operators []mitumbase.Address
}

func NewOperatorsBook(operators []mitumbase.Address) OperatorsBook {
	if operators == nil {
		return OperatorsBook{
			BaseHinter: hint.NewBaseHinter(OperatorsBookHint),
			operators:  []mitumbase.Address{},
		}
	}
	return OperatorsBook{
		BaseHinter: hint.NewBaseHinter(OperatorsBookHint),
		operators:  operators,
	}
}

func (ob OperatorsBook) IsValid([]byte) error {
	for i := range ob.operators {
		if err := ob.operators[i].IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (ob OperatorsBook) Bytes() []byte {
	ops := make([][]byte, len(ob.operators))

	for i, operator := range ob.operators {
		ops[i] = operator.Bytes()
	}

	return util.ConcatBytesSlice(ops...)
}

func (ob OperatorsBook) Hash() util.Hash {
	return ob.GenerateHash()
}

func (ob OperatorsBook) GenerateHash() util.Hash {
	return valuehash.NewSHA256(ob.Bytes())
}

func (ob OperatorsBook) IsEmpty() bool {
	return len(ob.operators) < 1
}

func (ob OperatorsBook) Equal(b OperatorsBook) bool {
	ob.Sort(true)
	b.Sort(true)

	for i := range ob.operators {
		if !ob.operators[i].Equal(b.operators[i]) {
			return false
		}
	}

	return true
}

func (ob *OperatorsBook) Sort(ascending bool) {
	sort.Slice(ob.operators, func(i, j int) bool {
		if ascending {
			return bytes.Compare(ob.operators[j].Bytes(), ob.operators[i].Bytes()) > 0
		}

		return bytes.Compare(ob.operators[j].Bytes(), ob.operators[i].Bytes()) < 0
	})
}

func (ob OperatorsBook) Exists(ag mitumbase.Address) bool {
	if ob.IsEmpty() {
		return false
	}

	for _, operator := range ob.operators {
		if ag.Equal(operator) {
			return true
		}
	}

	return false
}

func (ob OperatorsBook) Get(ag mitumbase.Address) (mitumbase.Address, error) {
	for _, operator := range ob.operators {
		if ag.Equal(operator) {
			return operator, nil
		}
	}

	return currencytypes.Address{}, errors.Errorf("account not in operators book, %q", ag)
}

func (ob *OperatorsBook) Append(ag mitumbase.Address) error {
	if err := ag.IsValid(nil); err != nil {
		return err
	}

	if ob.Exists(ag) {
		return errors.Errorf("account already in operators book, %q", ag)
	}

	if len(ob.operators) >= MaxOperators {
		return errors.Errorf("max operators, %v", ag)
	}

	ob.operators = append(ob.operators, ag)

	return nil
}

func (ob *OperatorsBook) Remove(ag mitumbase.Address) error {
	if !ob.Exists(ag) {
		return errors.Errorf("account not in operators book, %q", ag)
	}

	for i := range ob.operators {
		if ag.String() == ob.operators[i].String() {
			ob.operators[i] = ob.operators[len(ob.operators)-1]
			ob.operators[len(ob.operators)-1] = currencytypes.Address{}
			ob.operators = ob.operators[:len(ob.operators)-1]

			return nil
		}
	}
	return nil
}

func (ob OperatorsBook) Operators() []mitumbase.Address {
	return ob.operators
}
