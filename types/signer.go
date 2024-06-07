package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var SignerHint = hint.MustNewHint("mitum-nft-signer-v0.0.1")

var MaxSignerShare uint = 100

type Signer struct {
	hint.BaseHinter
	address base.Address
	share   uint
	signed  bool
}

func NewSigner(account base.Address, share uint, signed bool) Signer {
	return Signer{
		BaseHinter: hint.NewBaseHinter(SignerHint),
		address:    account,
		share:      share,
		signed:     signed,
	}
}

func (sgn Signer) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		sgn.BaseHinter,
		sgn.address,
	); err != nil {
		return err
	}

	if sgn.share > MaxSignerShare {
		return util.ErrInvalid.Errorf("share over max, %d > %d", sgn.share, MaxSignerShare)
	}

	return nil
}

func (sgn Signer) Bytes() []byte {
	bs := []byte{}
	if sgn.signed {
		bs = append(bs, 1)
	} else {
		bs = append(bs, 0)
	}

	return util.ConcatBytesSlice(
		sgn.address.Bytes(),
		util.UintToBytes(sgn.share),
		bs,
	)
}

func (sgn Signer) Address() base.Address {
	return sgn.address
}

func (sgn Signer) Share() uint {
	return sgn.share
}

func (sgn Signer) Signed() bool {
	return sgn.signed
}

func (sgn Signer) Equal(csigner Signer) bool {
	if sgn.Share() != csigner.Share() {
		return false
	}

	if !sgn.Address().Equal(csigner.Address()) {
		return false
	}

	if sgn.Signed() != csigner.Signed() {
		return false
	}

	return true
}
