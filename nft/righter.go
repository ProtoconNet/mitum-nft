package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	RightHolerType   = hint.Type("mitum-nft-right-holder")
	RightHolerHint   = hint.NewHint(RightHolerType, "v0.0.1")
	RightHolerHinter = RightHoler{BaseHinter: hint.NewBaseHinter(RightHolerHint)}
)

type RightHoler struct {
	hint.BaseHinter
	account base.Address
	signed  bool
	clue    string
}

func NewRightHoler(account base.Address, signed bool, clue string) RightHoler {
	return RightHoler{
		BaseHinter: hint.NewBaseHinter(RightHolerHint),
		account:    account,
		signed:     signed,
		clue:       clue,
	}
}

func MustNewRightHoler(account base.Address, signed bool, clue string) RightHoler {
	r := NewRightHoler(account, signed, clue)

	if err := r.IsValid(nil); err != nil {
		panic(err)
	}

	return r
}

func (r RightHoler) Bytes() []byte {
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

func (r RightHoler) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, r.BaseHinter, r.account); err != nil {
		return err
	}

	return nil
}

func (r RightHoler) Account() base.Address {
	return r.account
}

func (r RightHoler) Signed() bool {
	return r.signed
}

func (r RightHoler) Clue() string {
	return r.clue
}
