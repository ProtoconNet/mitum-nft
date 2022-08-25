package collection

import (
	"strings"
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type testDelegate struct {
	suite.Suite
}

func (t *testDelegate) TestNew() {
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

func (t *testDelegate) TestDeplicateCollectionAgentPair() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	collection := extensioncurrency.ContractID("ABC")
	agent := MustAddress(util.UUID().String())

	items := []DelegateItem{
		NewDelegateItem(collection, agent, DelegateAllow, "MCC"),
		NewDelegateItem(collection, agent, DelegateCancel, "MCC"),
	}
	fact := NewDelegateFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = delegate.IsValid(nil)
	t.Contains(err.Error(), "duplicate collection-agent pair found")
}

func (t *testDelegate) TestDelegateCancel() {
	sender := MustAddress(util.UUID().String())
	agent := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []DelegateItem{
		NewDelegateItem(extensioncurrency.ContractID("ABC"), agent, DelegateCancel, "MCC"),
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

func (t *testDelegate) TestSameSenderAgent() {

}

func (t *testDelegate) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []DelegateItem{}
	fact := NewDelegateFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = delegate.IsValid(nil)
	t.Contains(err.Error(), "empty items for DelegateFact")
}

func (t *testDelegate) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	collection := extensioncurrency.ContractID("ABC")
	items := []DelegateItem{
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateAllow, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateCancel, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateCancel, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateCancel, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateCancel, "MCC"),
		NewDelegateItem(collection, MustAddress(util.UUID().String()), DelegateCancel, "MCC"),
	}
	fact := NewDelegateFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = delegate.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testDelegate) TestOverSizeMemo() {
	sender := MustAddress(util.UUID().String())
	agent := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []DelegateItem{NewDelegateItem(extensioncurrency.ContractID("ABC"), agent, DelegateAllow, "MCC")}
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

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	delegate, err := NewDelegate(fact, fs, memo)
	t.NoError(err)

	err = delegate.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestDelegates(t *testing.T) {
	suite.Run(t, new(testDelegate))
}
