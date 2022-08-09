package nft

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testSigner struct {
	suite.Suite
}

func (t *testSigner) newSigner(account base.Address, share uint, signed bool) Signer {
	return MustNewSigner(account, share, signed)
}

func (t *testSigner) TestNew() {
	signer := t.newSigner(NewTestAddress(), 50, false)

	t.Equal(signer.Share(), uint(50))
	t.NotNil(signer.Account())
	t.False(signer.Signed())
}

func (t *testSigner) TestShareOverMax() {
	signer := NewSigner(NewTestAddress(), MaxSignerShare+1, false)
	t.Equal(signer.Share(), MaxSignerShare+1)
	t.Error(signer.IsValid(nil))
}

func (t *testSigner) TestShareMax() {
	signer := NewSigner(NewTestAddress(), MaxSignerShare, false)
	t.Equal(signer.Share(), uint(MaxSignerShare))
	t.NoError(signer.IsValid(nil))
}

func (t *testSigner) TestShareZero() {
	signer := NewSigner(NewTestAddress(), 0, false)
	t.Equal(signer.Share(), uint(0))
	t.NoError(signer.IsValid(nil))
}

func (t *testSigner) TestSigned() {
	signerSigned := t.newSigner(NewTestAddress(), 50, true)
	signerUnsigned := t.newSigner(NewTestAddress(), 50, false)

	t.True(signerSigned.Signed())
	t.False(signerUnsigned.Signed())
}

func (t *testSigner) TestEqual() {
	account := NewTestAddress()
	share := uint(10)
	signed := false

	signer0 := t.newSigner(account, share, signed)
	signer1 := t.newSigner(account, share, signed)
	t.True(signer0.Equal(signer1))

	signer2 := t.newSigner(NewTestAddress(), share, signed)
	t.False(signer0.Equal(signer2))

	signer3 := t.newSigner(account, 1, signed)
	t.False(signer0.Equal(signer3))

	signer4 := t.newSigner(account, share, true)
	t.False(signer0.Equal(signer4))
}

func TestSigner(t *testing.T) {
	suite.Run(t, new(testSigner))
}

type testSignerEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func (t *testSignerEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(currency.AddressHinter)
	encs.TestAddHinter(SignerHinter)
}

func (t *testSignerEncode) TestMarshal() {
	signer := NewSigner(NewTestAddress(), 50, false)
	t.NoError(signer.IsValid(nil))

	b, err := t.enc.Marshal(signer)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	usigner, ok := hinter.(Signer)
	t.True(ok)

	t.True(signer.Account().Equal(usigner.Account()))
	t.Equal(signer.Share(), usigner.Share())
	t.Equal(signer.Signed(), usigner.Signed())
}

func TestSignerEncodeJSON(t *testing.T) {
	b := new(testSignerEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestSignerEncodeBSON(t *testing.T) {
	b := new(testSignerEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
