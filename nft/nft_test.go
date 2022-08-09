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

type testNFT struct {
	suite.Suite
}

func (t *testNFT) newNFT(
	id NFTID,
	active bool,
	owner base.Address,
	hash NFTHash,
	uri URI,
	approved base.Address,
	creators Signers,
	copyrighters Signers,
) NFT {
	return MustNewNFT(id, active, owner, hash, uri, approved, creators, copyrighters)
}

func (t *testNFT) TestNew() {
	n := t.newNFT(
		NewTestNFTID(1),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)

	t.NoError(n.IsValid(nil))
	t.True(n.Active())
	t.NotNil(n.Hash())
	t.NotNil(n.Owner())
	t.NotNil(n.Approved())
	t.NotNil(n.Creators().Signers())
	t.NotNil(n.Copyrighters().Signers())
}

func (t *testNFT) TestExistsApproved() {
	owner := NewTestAddress()

	unapprovedNFT := t.newNFT(
		NewTestNFTID(1),
		true,
		owner,
		NFTHash(NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft",
		owner,
		NewTestSigners(),
		NewTestSigners(),
	)
	approvedNFT := t.newNFT(
		NewTestNFTID(2),
		true,
		owner,
		NFTHash(NewTestNFTID(2).Hash().String()),
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)

	t.NotNil(unapprovedNFT.Hash())
	t.NotNil(approvedNFT.Hash())

	t.False(unapprovedNFT.ExistsApproved())
	t.True(approvedNFT.ExistsApproved())
}

func (t *testNFT) TestActive() {
	activeNFT := t.newNFT(
		NewTestNFTID(1),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	deactiveNFT := t.newNFT(
		NewTestNFTID(2),
		false,
		NewTestAddress(),
		NFTHash(NewTestNFTID(2).Hash().String()),
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)

	t.NotNil(activeNFT.Hash())
	t.NotNil(deactiveNFT.Hash())

	t.True(activeNFT.Active())
	t.False(deactiveNFT.Active())

	t.NotEqual(activeNFT.Active(), deactiveNFT.Active())
}

func (t *testNFT) TestNFTHash() {
	noHashNFT := NewNFT(
		NewTestNFTID(1),
		true,
		NewTestAddress(),
		"",
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.NoError(noHashNFT.IsValid(nil))
	t.NotNil(noHashNFT.Hash())

	spaceHashNFT := NewNFT(
		NewTestNFTID(2),
		true,
		NewTestAddress(),
		"      ",
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.Error(spaceHashNFT.IsValid(nil))

	notTrimmedHashNFT := NewNFT(
		NewTestNFTID(3),
		true,
		NewTestAddress(),
		"     not trimmed nft hash   ",
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.NoError(notTrimmedHashNFT.IsValid(nil))
	t.NotNil(notTrimmedHashNFT.Hash())
}

func (t *testNFT) TestURI() {
	noUriNFT := NewNFT(
		NewTestNFTID(1),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(1).Hash().String()),
		"",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.Error(noUriNFT.IsValid(nil))

	spaceUriNFT := NewNFT(
		NewTestNFTID(2),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(2).Hash().String()),
		"     ",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.Error(spaceUriNFT.IsValid(nil))

	notTrimmedUriNFT := NewNFT(
		NewTestNFTID(3),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(3).Hash().String()),
		"      https://localhost:5000/nft       ",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.Error(notTrimmedUriNFT.IsValid(nil))

	notUriNFT := NewNFT(
		NewTestNFTID(4),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(4).Hash().String()),
		"abcdefg!@#$%^&*()-=[]",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.Error(notUriNFT.IsValid(nil))
}

func (t *testNFT) TestEqual() {
	nid := NewTestNFTID(1)
	active := true
	owner := NewTestAddress()
	hash := NFTHash(nid.Hash().String())
	uri := URI("https://localhost:5000/nft")
	approved := NewTestAddress()
	creators := NewTestSigners()
	copyrighters := NewTestSigners()

	n0 := t.newNFT(nid, active, owner, hash, uri, approved, creators, copyrighters)
	n1 := t.newNFT(nid, active, owner, hash, uri, approved, creators, copyrighters)
	t.NotNil(n0.Hash())
	t.NotNil(n1.Hash())
	t.True(n0.Hash().Equal(n1.Hash()))
	t.True(n0.Equal(n1))

	n2 := t.newNFT(NewTestNFTID(2), active, owner, hash, uri, approved, creators, copyrighters)
	t.NotNil(n2.Hash())
	t.False(n0.Hash().Equal(n2.Hash()))
	t.False(n0.Equal(n2))

	n3 := t.newNFT(nid, false, owner, hash, uri, approved, creators, copyrighters)
	t.NotNil(n3.Hash())
	t.False(n0.Hash().Equal(n3.Hash()))
	t.False(n0.Equal(n3))

	n4 := t.newNFT(nid, active, NewTestAddress(), hash, uri, approved, creators, copyrighters)
	t.NotNil(n4.Hash())
	t.False(n0.Hash().Equal(n4.Hash()))
	t.False(n0.Equal(n4))

	n5 := t.newNFT(nid, active, owner, NFTHash(NewTestNFTID(5).Hash().String()), uri, approved, creators, copyrighters)
	t.NotNil(n5.Hash())
	t.False(n0.Hash().Equal(n5.Hash()))
	t.False(n0.Equal(n5))

	n6 := t.newNFT(nid, active, owner, hash, "https://localhost:4000/nft", approved, creators, copyrighters)
	t.NotNil(n6.Hash())
	t.False(n0.Hash().Equal(n6.Hash()))
	t.False(n0.Equal(n6))

	n7 := t.newNFT(nid, active, owner, hash, uri, approved, NewTestSigners(), copyrighters)
	t.NotNil(n7.Hash())
	t.False(n0.Hash().Equal(n7.Hash()))
	t.False(n0.Equal(n7))

	n8 := t.newNFT(nid, active, owner, hash, uri, approved, creators, NewTestSigners())
	t.NotNil(n8.Hash())
	t.False(n0.Hash().Equal(n8.Hash()))
	t.False(n0.Equal(n8))
}

func TestNFT(t *testing.T) {
	suite.Run(t, new(testNFT))
}

type testNFTEncode struct {
	suite.Suite
	enc encoder.Encoder
}

func (t *testNFTEncode) SetupSuite() {
	encs := encoder.NewEncoders()
	encs.AddEncoder(t.enc)

	encs.TestAddHinter(currency.AddressHinter)
	encs.TestAddHinter(SignerHinter)
	encs.TestAddHinter(SignersHinter)
	encs.TestAddHinter(NFTIDHinter)
	encs.TestAddHinter(NFTHinter)
}

func (t *testNFTEncode) TestMarshal() {
	n := NewNFT(
		NewTestNFTID(1),
		true,
		NewTestAddress(),
		NFTHash(NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft",
		NewTestAddress(),
		NewTestSigners(),
		NewTestSigners(),
	)
	t.NoError(n.IsValid(nil))
	t.NotNil(n.Hash())

	b, err := t.enc.Marshal(n)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)
	un, ok := hinter.(NFT)
	t.True(ok)

	t.NotNil(un.Hash())

	t.True(n.Equal(un))
	t.True(n.Hash().Equal(un.Hash()))

	t.True(n.ID().Equal(un.ID()))
	t.True(n.Owner().Equal(un.Owner()))
	t.True(n.Approved().Equal(un.Approved()))
	t.True(n.Creators().Equal(un.Creators()))
	t.True(n.Copyrighters().Equal(un.Copyrighters()))
	t.Equal(n.Active(), un.Active())
	t.Equal(n.NftHash(), un.NftHash())
	t.Equal(n.Uri(), un.Uri())
}

func TestNFTEncodeJSON(t *testing.T) {
	b := new(testNFTEncode)
	b.enc = jsonenc.NewEncoder()

	suite.Run(t, b)
}

func TestNFTEncodeBSON(t *testing.T) {
	b := new(testNFTEncode)
	b.enc = bsonenc.NewEncoder()

	suite.Run(t, b)
}
