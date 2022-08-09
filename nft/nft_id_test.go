package nft

import (
	"math"
	"strings"
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testNFTID struct {
	suite.Suite
}

func (t *testNFTID) newNFTID(collection extensioncurrency.ContractID, idx uint64) NFTID {
	return MustNewNFTID(collection, idx)
}

func (t *testNFTID) TestNew() {
	collection := extensioncurrency.ContractID("ABC")
	idx := uint64(5000)

	nid := t.newNFTID(collection, idx)
	t.NoError(nid.IsValid(nil))
	t.Greater(nid.idx, uint64(0))
	t.NotNil(nid.Hash())
}

func (t *testNFTID) TestShortCollection() {
	collection := extensioncurrency.ContractID("AB")

	nid := NewNFTID(collection, 1)
	t.Equal(nid.Collection(), collection)
	t.NotNil(nid.Hash())
	t.Error(nid.IsValid(nil))
}

func (t *testNFTID) TestOverMaxCollection() {
	collection := extensioncurrency.ContractID(strings.Repeat("A", extensioncurrency.MaxLengthContractID+1))

	nid := NewNFTID(collection, 1)
	t.Equal(nid.Collection(), collection)
	t.NotNil(nid.Hash())
	t.Error(nid.IsValid(nil))
}

func (t *testNFTID) TestIDXOverZero() {
	zero := uint64(0)
	nidZeroIdx := NewNFTID(extensioncurrency.ContractID("ABC"), 0)
	nidPositiveIdx := NewNFTID(extensioncurrency.ContractID("ABC"), 10)

	t.Equal(nidZeroIdx.Idx(), zero)
	t.Greater(nidPositiveIdx.Idx(), zero)

	t.Error(nidZeroIdx.IsValid(nil))
	t.NoError(nidPositiveIdx.IsValid(nil))
	t.NotNil(nidPositiveIdx.Hash())
}

func (t *testNFTID) TestIDXOverMax() {
	nidMaxIdx := NewNFTID(extensioncurrency.ContractID("ABC"), uint64(MaxNFTIdx))
	t.Equal(uint64(MaxNFTIdx), nidMaxIdx.Idx())
	t.NoError(nidMaxIdx.IsValid(nil))
	t.NotNil(nidMaxIdx.Hash())

	if uint64(MaxNFTIdx) < math.MaxUint64 {
		nidOverMaxIdx := NewNFTID(extensioncurrency.ContractID("ABC"), uint64(MaxNFTIdx)+1)
		t.Equal(nidOverMaxIdx.Idx(), uint64(MaxNFTIdx)+1)
		t.Error(nidOverMaxIdx.IsValid(nil))
	}
}

func (t *testNFTID) TestEqual() {
	collection := extensioncurrency.ContractID("ABC")
	idx := uint64(1)
	zero := uint64(0)

	nid0 := t.newNFTID(collection, idx)
	t.Greater(nid0.idx, zero)
	t.NotNil(nid0.Hash())

	nid1 := t.newNFTID(collection, idx)
	t.Greater(nid1.idx, zero)
	t.NotNil(nid1.Hash())

	t.True(nid0.Equal(nid1))

	// different collection
	nid2 := t.newNFTID(extensioncurrency.ContractID("AAA"), idx)
	t.Greater(nid2.idx, zero)
	t.NotNil(nid2.Hash())

	t.False(nid0.Equal(nid2))

	// different idx
	nid3 := t.newNFTID(collection, 2)
	t.Greater(nid3.idx, zero)
	t.NotNil(nid3.Hash())

	t.False(nid0.Equal(nid3))
}

func TestNFTID(t *testing.T) {
	suite.Run(t, new(testNFTID))
}

type testNFTIDEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func (t *testNFTIDEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(NFTIDHinter)
}

func (t *testNFTIDEncode) TestMarshal() {
	nid := NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	t.NoError(nid.IsValid(nil))
	t.NotNil(nid.Hash())

	b, err := t.enc.Marshal(nid)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	unid, ok := hinter.(NFTID)
	t.True(ok)

	t.NotNil(unid.Hash())

	t.True(nid.Equal(unid))
	t.True(nid.Hash().Equal(unid.Hash()))

	t.Equal(nid.Collection(), unid.Collection())
	t.Equal(nid.Idx(), unid.Idx())
}

func TestNFTIDEncodeJSON(t *testing.T) {
	b := new(testNFTIDEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestNFTIDEncodeBSON(t *testing.T) {
	b := new(testNFTIDEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
