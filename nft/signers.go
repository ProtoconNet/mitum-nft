package nft

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MaxTotalShare uint = 100
	MaxSigners         = 10
)

var (
	SignersType   = hint.Type("mitum-nft-signers")
	SignersHint   = hint.NewHint(SignersType, "v0.0.1")
	SignersHinter = Signers{BaseHinter: hint.NewBaseHinter(SignersHint)}
)

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

func MustNewSigners(total uint, signers []Signer) Signers {
	sns := NewSigners(total, signers)

	if err := sns.IsValid(nil); err != nil {
		panic(err)
	}

	return sns
}

func (signers Signers) Bytes() []byte {
	bs := make([][]byte, len(signers.signers))

	for i := range signers.signers {
		bs[i] = signers.signers[i].Bytes()
	}

	return util.ConcatBytesSlice(
		util.UintToBytes(signers.total),
		util.ConcatBytesSlice(bs...),
	)
}

func (signers Signers) IsValid([]byte) error {
	if err := signers.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if signers.total > MaxTotalShare {
		return isvalid.InvalidError.Errorf("total share is over max; %d > %d", signers.total, MaxTotalShare)
	}

	if l := len(signers.signers); l > MaxSigners {
		return isvalid.InvalidError.Errorf("signers over allowed; %d > %d", l, MaxSigners)
	}

	var total uint = 0
	founds := map[base.Address]struct{}{}
	for i := range signers.signers {
		if err := signers.signers[i].IsValid(nil); err != nil {
			return err
		}

		acc := signers.signers[i].Account()
		if _, found := founds[acc]; found {
			return isvalid.InvalidError.Errorf("duplicate signer found; %q", acc)
		}
		founds[acc] = struct{}{}

		total += signers.signers[i].Share()
	}

	if total != signers.total {
		return isvalid.InvalidError.Errorf("total share must be equal to the sum of all shares; %d != %d", signers.total, total)
	}

	return nil
}

func (signers Signers) Total() uint {
	return signers.total
}

func (signers Signers) Signers() []Signer {
	return signers.signers
}

func (signers Signers) Addresses() []base.Address {
	as := make([]base.Address, len(signers.signers))
	for i := range signers.signers {
		as[i] = signers.signers[i].Account()
	}
	return as
}

func (signers Signers) Index(signer Signer) int {
	return signers.IndexByAddress(signer.Account())
}

func (signers Signers) IndexByAddress(address base.Address) int {
	for i := range signers.signers {
		if address.Equal(signers.signers[i].Account()) {
			return i
		}
	}
	return -1
}

func (signers Signers) Exists(signer Signer) bool {
	if idx := signers.Index(signer); idx >= 0 {
		return true
	}
	return false
}

func (signers Signers) IsSigned(signer Signer) bool {
	return signers.IsSignedByAddress(signer.Account())
}

func (signers Signers) IsSignedByAddress(address base.Address) bool {
	idx := signers.IndexByAddress(address)
	if idx < 0 {
		return false
	}
	return signers.signers[idx].Signed()
}

func (signers *Signers) SetSigner(signer Signer) error {
	idx := signers.Index(signer)
	if idx < 0 {
		return errors.Errorf("signer doesn't exist; %q", signer.Account())
	}
	signers.signers[idx] = signer
	return nil
}
