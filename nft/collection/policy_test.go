package collection

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testCollectionPolicy struct {
	suite.Suite
}

func (t *testCollectionPolicy) newCollectionPolicy(name CollectionName, royalty nft.PaymentParameter, uri nft.URI, whites []base.Address) CollectionPolicy {
	return MustNewCollectionPolicy(name, royalty, uri, whites)
}

func (t *testCollectionPolicy) TestNew() {
	policy := t.newCollectionPolicy("Collection", 0, "https://localhost:5000/collection", []base.Address{nft.NewTestAddress()})
	t.NotEmpty(policy.Name())
	t.LessOrEqual(policy.Royalty(), nft.MaxPaymentParameter)
	t.NotNil(policy.Whites())

	addresses, err := policy.Addresses()
	t.Nil(err)
	t.True(len(addresses) == 1)
}

func (t *testCollectionPolicy) TestShortName() {
	policy := NewCollectionPolicy("Co", 0, "https://localhost:5000/collection", []base.Address{nft.NewTestAddress()})
	t.True(len(policy.Name()) == 2)
	t.Error(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestOverMaxName() {
	name := strings.Repeat("a", MaxLengthCollectionName+1)
	policy := NewCollectionPolicy(CollectionName(name), 0, "https://localhost:5000/collection", []base.Address{nft.NewTestAddress()})
	t.True(len(policy.name) == MaxLengthCollectionName+1)
	t.Error(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestEmptyUri() {
	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{nft.NewTestAddress()})
	t.Empty(policy.Uri())
	t.NoError(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestOverMaxUri() {
	uri := strings.Repeat("a", nft.MaxURILength+1)
	policy := NewCollectionPolicy("Collection", 0, nft.URI(uri), []base.Address{nft.NewTestAddress()})
	t.True(len(policy.Uri()) == nft.MaxURILength+1)
	t.Error(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestOverMaxRoyalty() {
	policy := NewCollectionPolicy("Collection", nft.PaymentParameter(nft.MaxPaymentParameter+1), "https://localhost:5000/collection", []base.Address{nft.NewTestAddress()})
	t.True(policy.Royalty() == nft.PaymentParameter(nft.MaxPaymentParameter+1))
	t.Error(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestEmptyWhites() {
	policy := NewCollectionPolicy("Collection", 0, "https://localhost:5000/collection", []base.Address{})
	t.NotNil(policy.Whites())
	t.Empty(policy.Whites())
	t.True(len(policy.Whites()) == 0)

	addresses, err := policy.Addresses()
	t.Nil(err)
	t.True(len(addresses) == 0)
	t.NoError(policy.IsValid(nil))

}

func (t *testCollectionPolicy) TestOverMaxWhites() {
	policy := NewCollectionPolicy("Collection", 0, "https://localhost:5000/collection", []base.Address{
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
		nft.NewTestAddress(),
	})
	t.True(len(policy.Whites()) == MaxWhiteAddress+1)
	t.Error(policy.IsValid(nil))
}

func (t *testCollectionPolicy) TestEqual() {
	name := CollectionName("Collection")
	royalty := nft.PaymentParameter(0)
	uri := nft.URI("https://localhost:5000/collection")
	whites := []base.Address{nft.NewTestAddress()}

	p1 := t.newCollectionPolicy(name, royalty, uri, whites)
	p2 := t.newCollectionPolicy(name, royalty, uri, whites)
	t.True(p1.Equal(p2))

	// different name
	p3 := t.newCollectionPolicy("Different Collection", royalty, uri, whites)
	t.False(p1.Equal(p3))

	// different royalty
	p4 := t.newCollectionPolicy(name, nft.PaymentParameter(1), uri, whites)
	t.False(p1.Equal(p4))

	// different uri
	p5 := t.newCollectionPolicy(name, royalty, "", whites)
	t.False(p1.Equal(p5))

	// different whites
	p6 := t.newCollectionPolicy(name, royalty, uri, []base.Address{nft.NewTestAddress()})
	t.False(p1.Equal(p6))
}

type testCollectionPolicyEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func TestCollectionPolicy(t *testing.T) {
	suite.Run(t, new(testCollectionPolicy))
}

func (t *testCollectionPolicyEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(currency.AddressHinter)
	encs.TestAddHinter(CollectionPolicyHinter)
}

func (t *testCollectionPolicyEncode) TestMarshal() {
	policy := NewCollectionPolicy("Collection", 0, "https://localhost:5000/collection", []base.Address{nft.NewTestAddress()})
	t.NoError(policy.IsValid(nil))

	b, err := t.enc.Marshal(policy)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	upolicy, ok := hinter.(CollectionPolicy)
	t.True(ok)

	t.True(policy.Equal(upolicy))

	t.True(policy.Name() == upolicy.Name())
	t.True(policy.Royalty() == upolicy.Royalty())
	t.True(policy.Uri() == upolicy.Uri())

	whites := policy.Whites()
	uwhites := upolicy.Whites()

	t.Equal(len(whites), len(uwhites))

	sort.Slice(whites, func(i, j int) bool {
		return bytes.Compare(whites[j].Bytes(), whites[i].Bytes()) < 0
	})
	sort.Slice(uwhites, func(i, j int) bool {
		return bytes.Compare(uwhites[j].Bytes(), uwhites[i].Bytes()) < 0
	})

	for i := range whites {
		t.True(whites[i].Equal(uwhites[i]))
	}
}

func TestCollectionPolicyEncodeJSON(t *testing.T) {
	b := new(testCollectionPolicyEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestCollectionPolicyEncodeBSON(t *testing.T) {
	b := new(testCollectionPolicyEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
