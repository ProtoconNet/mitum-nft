package collection

import (
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testDelegateItem struct {
	suite.Suite
}

func (t *testDelegateItem) TestNew() {
	sender := MustAddress(util.UUID().String())
	agent := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []DelegateItem{
		NewDelegateItem(extensioncurrency.ContractID("ABC"), agent, DelegateAllow, "MCC"),
	}
	fact := NewDelegateFact(token, sender, items)

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

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	t.NoError(delegate.IsValid(nil))

	t.Implements((*base.Fact)(nil), delegate.Fact())
	t.Implements((*operation.Operation)(nil), delegate)
}

func (t *testDelegateItem) TestWrongMode() {
	mode := DelegateMode("wrong")

	sender := MustAddress(util.UUID().String())
	agent := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []DelegateItem{
		NewDelegateItem(extensioncurrency.ContractID("ABC"), agent, mode, "MCC"),
	}
	fact := NewDelegateFact(token, sender, items)

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

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = delegate.IsValid(nil)
	t.Contains(err.Error(), "wrong delegate mode")
}

func TestDelegateItem(t *testing.T) {
	suite.Run(t, new(testDelegateItem))
}

func testDelegateItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sender := MustAddress(util.UUID().String())
		agent := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		items := []DelegateItem{
			NewDelegateItem(extensioncurrency.ContractID("ABC"), agent, DelegateAllow, "MCC"),
		}
		fact := NewDelegateFact(token, sender, items)

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

		delegate, err := NewDelegate(fact, fs, "")
		t.NoError(err)

		return delegate
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Delegate)
		tb := b.(Delegate)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(DelegateFact)
		ufact := tb.Fact().(DelegateFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.Collection() == b.Collection())
			t.True(a.Agent().Equal(b.Agent()))
			t.True(a.Mode().Equal(b.Mode()))
			t.True(a.Currency() == b.Currency())
		}
	}

	return t
}

func TestDelegateItemEncodeJSON(t *testing.T) {
	suite.Run(t, testDelegateItemEncode(jsonenc.NewEncoder()))
}

func TestDelegateItemEncodeBSON(t *testing.T) {
	suite.Run(t, testDelegateItemEncode(bsonenc.NewEncoder()))
}
