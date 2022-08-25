package collection

import (
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testSignItem struct {
	suite.Suite
}

// test creator qualification
func (t *testSignItem) TestNew() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []SignItem{NewSignItem(CreatorQualification, nid, "MCC")}
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	t.NoError(sign.IsValid(nil))

	t.Implements((*base.Fact)(nil), sign.Fact())
	t.Implements((*operation.Operation)(nil), sign)
}

func (t *testSignItem) TestCopyrighterQualification() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []SignItem{NewSignItem(CopyrighterQualification, nid, "MCC")}
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	t.NoError(sign.IsValid(nil))

	t.Implements((*base.Fact)(nil), sign.Fact())
	t.Implements((*operation.Operation)(nil), sign)
}

func (t *testSignItem) TestWrongQualification() {
	qualification := Qualification("wrong")

	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []SignItem{NewSignItem(qualification, nid, "MCC")}
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "invalid qualification")
}

func (t *testSignItem) TestZeroIDX() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 0)
	items := []SignItem{NewSignItem(CreatorQualification, nid, "MCC")}
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "nid idx must be over zero")
}

func (t *testSignItem) TestOverMaxIDX() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), uint64(nft.MaxNFTIdx)+1)
	items := []SignItem{NewSignItem(CreatorQualification, nid, "MCC")}
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "nid idx over max")
}

func TestSignItem(t *testing.T) {
	suite.Run(t, new(testSignItem))
}

func testSignItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sender := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		nid0 := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
		nid1 := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 2)
		items := []SignItem{
			NewSignItem(CreatorQualification, nid0, "MCC"),
			NewSignItem(CopyrighterQualification, nid1, "MCC"),
		}
		fact := NewSignFact(token, sender, items)

		var fs []base.FactSign

		for _, pk := range []key.Privatekey{
			key.NewBasePrivatekey(),
			key.NewBasePrivatekey(),
			key.NewBasePrivatekey(),
		} {
			sig, err := base.NewFactSignature(pk, fact, nil)
			t.NoError(err)

			fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
		}

		sign, err := NewSign(fact, fs, "")
		t.NoError(err)

		return sign
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Sign)
		tb := b.(Sign)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(SignFact)
		ufact := tb.Fact().(SignFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.NFT().Equal(b.NFT()))
			t.Equal(a.Qualification(), b.Qualification())
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestSignItemEncodeJSON(t *testing.T) {
	suite.Run(t, testSignItemEncode(jsonenc.NewEncoder()))
}

func TestSignItemEncodeBSON(t *testing.T) {
	suite.Run(t, testSignItemEncode(bsonenc.NewEncoder()))
}
