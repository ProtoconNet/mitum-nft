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

type testBurnItem struct {
	suite.Suite
}

func (t *testBurnItem) TestNew() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []BurnItem{NewBurnItem(nid, "MCC")}
	fact := NewBurnFact(token, sender, items)

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

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	t.NoError(burn.IsValid(nil))

	t.Implements((*base.Fact)(nil), burn.Fact())
	t.Implements((*operation.Operation)(nil), burn)
}

func (t *testBurnItem) TestZeroIDX() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 0)
	items := []BurnItem{NewBurnItem(nid, "MCC")}
	fact := NewBurnFact(token, sender, items)

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

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	err = burn.IsValid(nil)
	t.Contains(err.Error(), "nid idx must be over zero")
}

func (t *testBurnItem) TestOverMaxIDX() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), nft.MaxNFTIdx+1)
	items := []BurnItem{NewBurnItem(nid, "MCC")}
	fact := NewBurnFact(token, sender, items)

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

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	err = burn.IsValid(nil)
	t.Contains(err.Error(), "nid idx over max")
}

func TestBurnItem(t *testing.T) {
	suite.Run(t, new(testBurnItem))
}

func testBurnItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sender := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
		items := []BurnItem{NewBurnItem(nid, "MCC")}
		fact := NewBurnFact(token, sender, items)

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

		burn, err := NewBurn(fact, fs, "")
		t.NoError(err)

		return burn
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Burn)
		tb := b.(Burn)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(BurnFact)
		ufact := tb.Fact().(BurnFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.NFT().Equal(b.NFT()))
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestBurnItemEncodeJSON(t *testing.T) {
	suite.Run(t, testBurnItemEncode(jsonenc.NewEncoder()))
}

func TestBurnItemEncodeBSON(t *testing.T) {
	suite.Run(t, testBurnItemEncode(bsonenc.NewEncoder()))
}
