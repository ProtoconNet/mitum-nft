package nft

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testSigners struct {
	suite.Suite
}

func (t *testSigners) newSigners(total uint, signers []Signer) Signers {
	return MustNewSigners(total, signers)
}

func (t *testSigners) TestNew() {
	signer0 := MustNewSigner(NewTestAddress(), 50, false)
	signer1 := MustNewSigner(NewTestAddress(), 50, false)

	signers := t.newSigners(100, []Signer{signer0, signer1})

	ads := signers.Addresses()
	t.Equal(len(signers.Signers()), len(ads))

	var total uint = 0
	for i := range signers.Signers() {
		total += signers.Signers()[i].Share()
	}
	t.Equal(signers.Total(), total)
}

func (t *testSigners) TestTotalOverMax() {
	signer0 := MustNewSigner(NewTestAddress(), 51, false)
	signer1 := MustNewSigner(NewTestAddress(), 50, false)
	signers := NewSigners(uint(MaxTotalShare+1), []Signer{signer0, signer1})

	t.Equal(signers.Total(), uint(MaxTotalShare+1))
	t.Error(signers.IsValid(nil))
}

func (t *testSigners) TestTotalMax() {
	signer0 := MustNewSigner(NewTestAddress(), 20, false)
	signer1 := MustNewSigner(NewTestAddress(), 80, false)
	signers := NewSigners(100, []Signer{signer0, signer1})

	t.Equal(signers.Total(), uint(100))
	t.NoError(signers.IsValid(nil))
}

func (t *testSigners) TestTotalDifferentShares() {
	signer0 := MustNewSigner(NewTestAddress(), 40, false)
	signer1 := MustNewSigner(NewTestAddress(), 40, false)
	signer2 := MustNewSigner(NewTestAddress(), 40, false)

	signers := NewSigners(120, []Signer{signer0, signer1, signer2})
	t.Equal(signers.Total(), uint(120))
	t.Error(signers.IsValid(nil))
}

func (t *testSigners) TestTotalZero() {
	signer := MustNewSigner(NewTestAddress(), 0, false)
	signers := NewSigners(0, []Signer{signer})
	t.Equal(signers.Total(), uint(0))
	t.NoError(signers.IsValid(nil))
}

func (t *testSigners) TestOverMaxSigners() {
	signers0 := make([]Signer, 10)
	for i := range signers0 {
		signers0[i] = MustNewSigner(NewTestAddress(), 10, false)
	}
	sgns0 := NewSigners(100, signers0)

	var total uint = 0
	for i := range sgns0.Signers() {
		total += sgns0.Signers()[i].Share()
	}
	t.Equal(len(sgns0.Signers()), 10)
	t.Equal(sgns0.Total(), uint(100))
	t.Equal(sgns0.Total(), total)
	t.NoError(sgns0.IsValid(nil))

	signers1 := make([]Signer, MaxSigners+1)
	for i := range signers1 {
		signers1[i] = MustNewSigner(NewTestAddress(), 0, false)
	}
	sgns1 := NewSigners(0, signers1)

	total = 0
	for i := range sgns1.Signers() {
		total += sgns1.Signers()[i].Share()
	}
	t.Equal(len(sgns1.Signers()), MaxSigners+1)
	t.Equal(sgns1.Total(), uint(0))
	t.Equal(sgns1.Total(), total)
	t.Error(sgns1.IsValid(nil))
}

func (t *testSigners) TestZeroSigners() {
	signers := NewSigners(0, []Signer{})
	t.Equal(signers.Total(), uint(0))
	t.NoError(signers.IsValid(nil))
}

func (t *testSigners) TestDuplicateSigner() {
	signer0 := MustNewSigner(NewTestAddress(), 10, false)
	signer1 := MustNewSigner(NewTestAddress(), 10, false)

	sgns := NewSigners(30, []Signer{signer0, signer0, signer1})

	var total uint = 0
	for i := range sgns.Signers() {
		total += sgns.Signers()[i].Share()
	}
	t.Equal(len(sgns.Signers()), 3)
	t.Equal(sgns.Total(), uint(30))
	t.Equal(sgns.Total(), total)
	t.Error(sgns.IsValid(nil))
}

func (t *testSigners) TestEqual() {
	total := uint(100)
	signers := []Signer{
		MustNewSigner(NewTestAddress(), 100, false),
	}

	signers0 := t.newSigners(total, signers)
	signers1 := t.newSigners(total, signers)
	t.True(signers0.Equal(signers1))

	signers2 := NewSigners(90, signers)
	t.Error(signers2.IsValid(nil))
	t.False(signers0.Equal(signers2))

	signers3 := t.newSigners(total, []Signer{
		MustNewSigner(NewTestAddress(), 50, false),
		MustNewSigner(NewTestAddress(), 50, false),
	})
	t.False(signers0.Equal(signers3))
}

func TestSigners(t *testing.T) {
	suite.Run(t, new(testSigners))
}

type testSignersEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func (t *testSignersEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(currency.AddressHinter)
	encs.TestAddHinter(SignerHinter)
	encs.TestAddHinter(SignersHinter)
}

func (t *testSignersEncode) TestMarshal() {
	signers := make([]Signer, 10)
	for i := range signers {
		signers[i] = MustNewSigner(NewTestAddress(), 10, false)
	}
	sgns := NewSigners(100, signers)
	t.NoError(sgns.IsValid(nil))

	b, err := t.enc.Marshal(sgns)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	usgns, ok := hinter.(Signers)
	t.True(ok)

	t.Equal(sgns.Total(), usgns.Total())
	t.Equal(len(sgns.Signers()), len(usgns.Signers()))

	sgnsAddresses := sgns.Addresses()
	usgnsAddresses := usgns.Addresses()

	t.Equal(len(sgnsAddresses), len(usgnsAddresses))

	var sgnsTotal uint = 0
	var usgnsTotal uint = 0

	for i := range sgns.Signers() {
		sgnsTotal += sgns.Signers()[i].Share()
	}

	for i := range usgns.Signers() {
		usgnsTotal += usgns.Signers()[i].Share()
	}

	t.Equal(sgnsTotal, usgnsTotal)
}

func TestSignersEncodeJSON(t *testing.T) {
	b := new(testSignersEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestSignersEncodeBSON(t *testing.T) {
	b := new(testSignersEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
