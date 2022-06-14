package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	RightHolderType   = hint.Type("mitum-nft-right-holder")
	RightHolderHint   = hint.NewHint(RightHolderType, "v0.0.1")
	RightHolderHinter = RightHolder{BaseHinter: hint.NewBaseHinter(RightHolderHint)}
)

type RightHolder struct {
	hint.BaseHinter
	account base.Address
	signed  bool
	clue    string
}

func NewRightHolder(account base.Address, signed bool, clue string) RightHolder {
	return RightHolder{
		BaseHinter: hint.NewBaseHinter(RightHolderHint),
		account:    account,
		signed:     signed,
		clue:       clue,
	}
}

func MustNewRightHolder(account base.Address, signed bool, clue string) RightHolder {
	r := NewRightHolder(account, signed, clue)

	if err := r.IsValid(nil); err != nil {
		panic(err)
	}

	return r
}

func (r RightHolder) Bytes() []byte {
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

func (r RightHolder) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, r.BaseHinter, r.account); err != nil {
		return err
	}

	return nil
}

func (r RightHolder) Account() base.Address {
	return r.account
}

func (r RightHolder) Signed() bool {
	return r.signed
}

func (r RightHolder) Clue() string {
	return r.clue
}
