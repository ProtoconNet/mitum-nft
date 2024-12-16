package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	MaxTotalShare uint = 100
	MaxSigners         = 10
)

var SignersHint = hint.MustNewHint("mitum-nft-signers-v0.0.1")

type Signers struct {
	hint.BaseHinter
	signers []Signer
}

func NewSigners(signers []Signer) Signers {
	return Signers{
		BaseHinter: hint.NewBaseHinter(SignersHint),
		signers:    signers,
	}
}

func (sgns Signers) IsValid([]byte) error {
	if err := sgns.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if l := len(sgns.signers); l > MaxSigners {
		return common.ErrValOOR.Wrap(errors.Errorf("signers over allowed, %d > %d", l, MaxSigners))
	}

	var total uint = 0
	founds := map[string]struct{}{}
	for _, signer := range sgns.signers {
		if err := signer.IsValid(nil); err != nil {
			return err
		}

		acc := signer.Address()
		if _, found := founds[acc.String()]; found {
			return common.ErrDupVal.Wrap(errors.Errorf("signer account %v", acc))
		}
		founds[acc.String()] = struct{}{}

		total += signer.Share()
	}

	if total > MaxTotalShare {
		return common.ErrValOOR.Wrap(errors.Errorf("signers total share over max, %d > %d", total, MaxTotalShare))
	}

	return nil
}

func (sgns Signers) Bytes() []byte {
	bs := make([][]byte, len(sgns.signers))

	for i, signer := range sgns.signers {
		bs[i] = signer.Bytes()
	}

	return util.ConcatBytesSlice(
		util.ConcatBytesSlice(bs...),
	)
}

func (sgns Signers) Signers() []Signer {
	return sgns.signers
}

func (sgns Signers) Addresses() []base.Address {
	as := make([]base.Address, len(sgns.signers))
	for i, signer := range sgns.signers {
		as[i] = signer.Address()
	}
	return as
}

func (sgns Signers) Index(signer Signer) int {
	return sgns.IndexByAddress(signer.Address())
}

func (sgns Signers) IndexByAddress(address base.Address) int {
	for i := range sgns.signers {
		if address.Equal(sgns.signers[i].Address()) {
			return i
		}
	}
	return -1
}

func (sgns Signers) Exists(signer Signer) bool {
	if idx := sgns.Index(signer); idx >= 0 {
		return true
	}
	return false
}

func (sgns Signers) Equal(ys Signers) bool {
	if len(sgns.Signers()) != len(ys.Signers()) {
		return false
	}

	xsg := sgns.Signers()
	sort.Slice(xsg, func(i, j int) bool {
		return bytes.Compare(xsg[j].Bytes(), xsg[i].Bytes()) < 0
	})

	ysg := ys.Signers()
	sort.Slice(ysg, func(i, j int) bool {
		return bytes.Compare(ysg[j].Bytes(), ysg[i].Bytes()) < 0
	})

	for i := range xsg {
		if !xsg[i].Equal(ysg[i]) {
			return false
		}
	}

	return true
}

func (sgns Signers) IsSigned(sgn Signer) bool {
	return sgns.IsSignedByAddress(sgn.Address())
}

func (sgns Signers) IsSignedByAddress(address base.Address) bool {
	idx := sgns.IndexByAddress(address)
	if idx < 0 {
		return false
	}
	return sgns.signers[idx].Signed()
}

func (sgns *Signers) SetSigner(sgn Signer) error {
	idx := sgns.Index(sgn)
	if idx < 0 {
		return errors.Errorf("signer not in signers, %v", sgn.Address())
	}
	sgns.signers[idx] = sgn
	return nil
}
