package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	SignerType   = hint.Type("mitum-nft-signer")
	SignerHint   = hint.NewHint(SignerType, "v0.0.1")
	SignerHinter = Signer{BaseHinter: hint.NewBaseHinter(SignerHint)}
)

type Signer struct {
	hint.BaseHinter
	account base.Address
	signed  bool
}

func NewSigner(account base.Address, signed bool) Signer {
	return Signer{
		BaseHinter: hint.NewBaseHinter(SignerHint),
		account:    account,
		signed:     signed,
	}
}

func MustNewSigner(account base.Address, signed bool, clue string) Signer {
	signer := NewSigner(account, signed)

	if err := signer.IsValid(nil); err != nil {
		panic(err)
	}

	return signer
}

func (signer Signer) Bytes() []byte {
	bs := []byte{}
	if signer.signed {
		bs = append(bs, 1)
	} else {
		bs = append(bs, 0)
	}

	return util.ConcatBytesSlice(
		signer.account.Bytes(),
		bs,
	)
}

func (signer Signer) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, signer.BaseHinter, signer.account); err != nil {
		return err
	}

	return nil
}

func (signer Signer) Account() base.Address {
	return signer.account
}

func (signer Signer) Signed() bool {
	return signer.signed
}
