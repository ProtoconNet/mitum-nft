package nft

import (
	"strings"
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testDesign struct {
	suite.Suite
}

func (t *testDesign) newDesign(parent base.Address, creator base.Address, symbol extensioncurrency.ContractID, active bool, policy BasePolicy) Design {
	return MustNewDesign(parent, creator, symbol, active, policy)
}

func (t *testDesign) TestNew() {
	design := t.newDesign(NewTestAddress(), NewTestAddress(), extensioncurrency.ContractID("ABC"), true, NewTestPolicy(1))
	t.NotNil(design.Parent())
	t.NotNil(design.Creator())
	t.NotNil(design.Hash())
	t.True(design.Active())

	addresses, err := design.Addresses()
	t.NoError(err)

	policyAddresses, err := design.Policy().Addresses()
	t.NoError(err)

	t.Equal(len(addresses), len(policyAddresses)+2)
}

func (t *testDesign) TestDesignActivation() {
	activeDesign := t.newDesign(NewTestAddress(), NewTestAddress(), extensioncurrency.ContractID("ABC"), true, NewTestPolicy(1))
	t.True(activeDesign.Active())
	t.NotNil(activeDesign.Hash())

	deactiveDesign := t.newDesign(NewTestAddress(), NewTestAddress(), extensioncurrency.ContractID("ABC"), false, NewTestPolicy(1))
	t.False(deactiveDesign.Active())
	t.NotNil(deactiveDesign.Hash())

	t.NotEqual(activeDesign.Active(), deactiveDesign.Active())
}

func (t *testDesign) TestSameParentCreator() {
	parent := NewTestAddress()

	design := NewDesign(parent, parent, extensioncurrency.ContractID("ABC"), true, NewTestPolicy(1))
	t.True(parent.Equal(design.Parent()))
	t.True(parent.Equal(design.Creator()))
	t.NotNil(design.Hash())

	t.Error(design.IsValid(nil))
}

func (t *testDesign) TestShortSymbol() {
	design := NewDesign(NewTestAddress(), NewTestAddress(), extensioncurrency.ContractID("AB"), true, NewTestPolicy(1))
	t.Equal(design.Symbol(), extensioncurrency.ContractID("AB"))
	t.NotNil(design.Hash())

	t.Error(design.IsValid(nil))
}

func (t *testDesign) TestOverMaxSymbol() {
	symbol := extensioncurrency.ContractID(strings.Repeat("A", extensioncurrency.MaxLengthContractID+1))
	design := NewDesign(NewTestAddress(), NewTestAddress(), symbol, true, NewTestPolicy(1))
	t.Equal(design.Symbol(), symbol)
	t.NotNil(design.Hash())

	t.Error(design.IsValid(nil))
}

func (t *testDesign) TestEqual() {
	collection := extensioncurrency.ContractID("ABC")
	parent := NewTestAddress()
	creator := NewTestAddress()
	policy := NewTestPolicy(1)

	d1 := t.newDesign(parent, creator, collection, true, policy)
	t.NotNil(d1.Hash())
	t.True(d1.Active())

	d2 := t.newDesign(parent, creator, collection, true, policy)
	t.NotNil(d2.Hash())
	t.True(d2.Active())

	t.True(d1.Equal(d2))

	// different parent
	d3 := t.newDesign(NewTestAddress(), creator, collection, true, policy)
	t.NotNil(d3.Hash())
	t.True(d3.Active())

	t.False(d1.Equal(d3))

	// different creator
	d4 := t.newDesign(parent, NewTestAddress(), collection, true, policy)
	t.NotNil(d4.Hash())
	t.True(d4.Active())

	t.False(d1.Equal(d4))

	// different collection
	d5 := t.newDesign(parent, creator, extensioncurrency.ContractID("AAA"), true, policy)
	t.NotNil(d5.Hash())
	t.True(d5.Active())

	t.False(d1.Equal(d5))

	// different active
	d6 := t.newDesign(parent, creator, collection, false, policy)
	t.NotNil(d6.Hash())
	t.False(d6.Active())

	t.False(d1.Equal(d6))

	// different policy
	d7 := t.newDesign(parent, creator, collection, true, NewTestPolicy(3))
	t.NotNil(d7.Hash())
	t.True(d7.Active())

	t.False(d1.Equal(d7))
}

type testDesignEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func TestDesign(t *testing.T) {
	suite.Run(t, new(testDesign))
}

func (t *testDesignEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(currency.AddressHinter)
	encs.TestAddHinter(DesignHinter)
	encs.TestAddHinter(TestPolicyHinter)
}

func (t *testDesignEncode) TestMarshal() {
	design := NewDesign(NewTestAddress(), NewTestAddress(), extensioncurrency.ContractID("ABC"), true, NewTestPolicy(1))
	t.NoError(design.IsValid(nil))
	t.NotNil(design.Hash())

	b, err := t.enc.Marshal(design)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	udesign, ok := hinter.(Design)
	t.True(ok)

	t.NotNil(udesign.Hash())

	t.True(design.Equal(udesign))
	t.True(design.Hash().Equal(udesign.Hash()))

	t.True(design.Parent().Equal(udesign.Parent()))
	t.True(design.Creator().Equal(udesign.Creator()))
	t.True(design.Policy().Equal(udesign.Policy()))
	t.Equal(design.Symbol(), udesign.Symbol())
	t.Equal(design.Active(), udesign.Active())
}

func TestDesignEncodeJSON(t *testing.T) {
	b := new(testDesignEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestDesignEncodeBSON(t *testing.T) {
	b := new(testDesignEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
