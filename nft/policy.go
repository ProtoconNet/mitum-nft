package nft

import (
	"regexp"

	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MinLengthSymbol = 3
	MaxLengthSymbol = 10
	ReValidSymbol   = regexp.MustCompile(`^[A-Z]+$`)
)

type Symbol string

func (s Symbol) Bytes() []byte {
	return []byte(s)
}

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) IsValid([]byte) error {
	if l := len(s); l < MinLengthSymbol || l > MaxLengthSymbol {
		return isvalid.InvalidError.Errorf(
			"invalid length of symbol; %d <= %d <= %d", MinLengthSymbol, l, MaxLengthSymbol)
	} else if !ReValidSymbol.Match([]byte(s)) {
		return isvalid.InvalidError.Errorf("wrong symbol; %q", s)
	}

	return nil
}

type PaymentParameter uint

func (pp PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(pp))
}

func (pp PaymentParameter) Uint() uint {
	return uint(pp)
}

func (pp PaymentParameter) IsValid([]byte) error {
	if uint(pp) >= 100 {
		return isvalid.InvalidError.Errorf(
			"invalid range of symbol; %d <= %d < %d", 0, pp, 100)
	}

	return nil
}

type BasePolicy interface {
	isvalid.IsValider
	Bytes() []byte
	Symbol() Symbol
	Equal(policy BasePolicy) bool
	Rebuild() BasePolicy
}
