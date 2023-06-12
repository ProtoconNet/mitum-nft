package digest

import (
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v3/types"
)

var (
	AccountValueHint = hint.MustNewHint("mitum-currency-account-value-v0.0.1")
)

type AccountValue struct {
	hint.BaseHinter
	ac      types.Account
	balance []types.Amount
	height  mitumbase.Height
}

func NewAccountValue(st mitumbase.State) (AccountValue, error) {
	var ac types.Account
	switch a, ok, err := IsAccountState(st); {
	case err != nil:
		return AccountValue{}, err
	case !ok:
		return AccountValue{}, errors.Errorf("not state for types.Account, %T", st.Value())
	default:
		ac = a
	}

	return AccountValue{
		BaseHinter: hint.NewBaseHinter(AccountValueHint),
		ac:         ac,
		height:     st.Height(),
	}, nil
}

func (va AccountValue) Account() types.Account {
	return va.ac
}

func (va AccountValue) Balance() []types.Amount {
	return va.balance
}

func (va AccountValue) Height() mitumbase.Height {
	return va.height
}

func (va AccountValue) SetHeight(height mitumbase.Height) AccountValue {
	va.height = height

	return va
}

func (va AccountValue) SetBalance(balance []types.Amount) AccountValue {
	va.balance = balance

	return va
}
