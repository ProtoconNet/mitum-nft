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

var MaxSignerShare uint = 100

type Signer struct {
	hint.BaseHinter
	account base.Address
	share   uint
	signed  bool
}

func NewSigner(account base.Address, share uint, signed bool) Signer {
	return Signer{
		BaseHinter: hint.NewBaseHinter(SignerHint),
		account:    account,
		share:      share,
		signed:     signed,
	}
}

func MustNewSigner(account base.Address, share uint, signed bool) Signer {
	signer := NewSigner(account, share, signed)

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
		util.UintToBytes(signer.share),
		bs,
	)
}

func (signer Signer) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, signer.BaseHinter, signer.account); err != nil {
		return err
	}

	if signer.share > MaxSignerShare {
		return isvalid.InvalidError.Errorf("share is over max; %d > %d", signer.share, MaxSignerShare)
	}

	return nil
}

func (signer Signer) Account() base.Address {
	return signer.account
}

func (signer Signer) Share() uint {
	return signer.share
}

func (signer Signer) Signed() bool {
	return signer.signed
}

func (signer Signer) Equal(csigner Signer) bool {
	if signer.Share() != csigner.Share() {
		return false
	}

	if !signer.Account().Equal(csigner.Account()) {
		return false
	}

	if signer.Signed() != csigner.Signed() {
		return false
	}

	return true
}
