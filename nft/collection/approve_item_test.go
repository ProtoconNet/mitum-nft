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

type testApproveItem struct {
	suite.Suite
}

func (t *testApproveItem) TestNew() {
	sender := MustAddress(util.UUID().String())
	approved := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []ApproveItem{NewApproveItem(approved, nid, "MCC")}
	fact := NewApproveFact(token, sender, items)

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

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	t.NoError(approve.IsValid(nil))

	t.Implements((*base.Fact)(nil), approve.Fact())
	t.Implements((*operation.Operation)(nil), approve)
}

func (t *testApproveItem) TestZeroIDX() {
	sender := MustAddress(util.UUID().String())
	approved := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 0)
	items := []ApproveItem{NewApproveItem(approved, nid, "MCC")}
	fact := NewApproveFact(token, sender, items)

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

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = approve.IsValid(nil)
	t.Contains(err.Error(), "nid idx must be over zero")
}

func (t *testApproveItem) TestOverMaxIDX() {
	sender := MustAddress(util.UUID().String())
	approved := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), nft.MaxNFTIdx+1)
	items := []ApproveItem{NewApproveItem(approved, nid, "MCC")}
	fact := NewApproveFact(token, sender, items)

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

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = approve.IsValid(nil)
	t.Contains(err.Error(), "nid idx over max")
}

func TestApproveItem(t *testing.T) {
	suite.Run(t, new(testApproveItem))
}

func testApproveItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sender := MustAddress(util.UUID().String())
		approved := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
		items := []ApproveItem{NewApproveItem(approved, nid, "MCC")}
		fact := NewApproveFact(token, sender, items)

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

		approve, err := NewApprove(fact, fs, "")
		t.NoError(err)

		return approve
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Approve)
		tb := b.(Approve)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(ApproveFact)
		ufact := tb.Fact().(ApproveFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.Approved().Equal(b.Approved()))
			t.True(a.NFT().Equal(b.NFT()))
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestApproveItemEncodeJSON(t *testing.T) {
	suite.Run(t, testApproveItemEncode(jsonenc.NewEncoder()))
}

func TestApproveItemEncodeBSON(t *testing.T) {
	suite.Run(t, testApproveItemEncode(bsonenc.NewEncoder()))
}
