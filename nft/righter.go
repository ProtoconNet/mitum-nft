package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	RighterType   = hint.Type("mitum-nft-righter")
	RighterHint   = hint.NewHint(RighterType, "v0.0.1")
	RighterHinter = Righter{BaseHinter: hint.NewBaseHinter(RighterHint)}
)

type Righter struct {
	hint.BaseHinter
	account base.Address
	signed  bool
	clue    string
}

func NewRighter(account base.Address, signed bool, clue string) Righter {
	return Righter{
		BaseHinter: hint.NewBaseHinter(RighterHint),
		account:    account,
		signed:     signed,
		clue:       clue,
	}
}

func MustNewRighter(account base.Address, signed bool, clue string) Righter {
	righter := NewRighter(account, signed, clue)

	if err := righter.IsValid(nil); err != nil {
		panic(err)
	}

	return righter
}

func (r Righter) Bytes() []byte {
	bs := []byte{}
	if r.signed {
		bs = append(bs, 1)
	} else {
		bs = append(bs, 0)
	}

	return util.ConcatBytesSlice(
		r.account.Bytes(),
		bs,
		[]byte(r.clue),
	)
}

func (r Righter) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, r.BaseHinter, r.account); err != nil {
		return err
	}

	return nil
}

func (r Righter) Account() base.Address {
	return r.account
}

func (r Righter) Signed() bool {
	return r.signed
}

func (r Righter) Clue() string {
	return r.clue
}
