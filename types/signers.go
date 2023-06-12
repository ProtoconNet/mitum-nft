package types

import (
	"bytes"
	"sort"

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
	total   uint
	signers []Signer
}

func NewSigners(total uint, signers []Signer) Signers {
	return Signers{
		BaseHinter: hint.NewBaseHinter(SignersHint),
		total:      total,
		signers:    signers,
	}
}

func (sgns Signers) IsValid([]byte) error {
	if err := sgns.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if sgns.total > MaxTotalShare {
		return util.ErrInvalid.Errorf("total share over max, %d > %d", sgns.total, MaxTotalShare)
	}

	if l := len(sgns.signers); l > MaxSigners {
		return util.ErrInvalid.Errorf("signers over allowed, %d > %d", l, MaxSigners)
	}

	var total uint = 0
	founds := map[string]struct{}{}
	for _, signer := range sgns.signers {
		if err := signer.IsValid(nil); err != nil {
			return err
		}

		acc := signer.Account()
		if _, found := founds[acc.String()]; found {
			return util.ErrInvalid.Errorf("duplicate signer found, %q", acc)
		}
		founds[acc.String()] = struct{}{}

		total += signer.Share()
	}

	if total != sgns.total {
		return util.ErrInvalid.Errorf("total share must be equal to the sum of all shares, %d != %d", sgns.total, total)
	}

	return nil
}

func (sgns Signers) Bytes() []byte {
	bs := make([][]byte, len(sgns.signers))

	for i, signer := range sgns.signers {
		bs[i] = signer.Bytes()
	}

	return util.ConcatBytesSlice(
		util.UintToBytes(sgns.total),
		util.ConcatBytesSlice(bs...),
	)
}

func (sgns Signers) Total() uint {
	return sgns.total
}

func (sgns Signers) Signers() []Signer {
	return sgns.signers
}

func (sgns Signers) Addresses() []base.Address {
	as := make([]base.Address, len(sgns.signers))
	for i, signer := range sgns.signers {
		as[i] = signer.Account()
	}
	return as
}

func (sgns Signers) Index(signer Signer) int {
	return sgns.IndexByAddress(signer.Account())
}

func (sgns Signers) IndexByAddress(address base.Address) int {
	for i := range sgns.signers {
		if address.Equal(sgns.signers[i].Account()) {
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

func (xs Signers) Equal(ys Signers) bool {
	if xs.Total() != ys.Total() {
		return false
	}

	if len(xs.Signers()) != len(ys.Signers()) {
		return false
	}

	xsg := xs.Signers()
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
	return sgns.IsSignedByAddress(sgn.Account())
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
		return errors.Errorf("signer not in signers, %q", sgn.Account())
	}
	sgns.signers[idx] = sgn
	return nil
}
