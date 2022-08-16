package collection

import (
	"net/url"
	"strings"
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type testCollectionPolicyUpdater struct {
	suite.Suite
}

func (t *testCollectionPolicyUpdater) TestNew() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "https://localhost:5000/collection", []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionPolicyUpdater.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionPolicyUpdater.Fact())
	t.Implements((*operation.Operation)(nil), collectionPolicyUpdater)
}

func (t *testCollectionPolicyUpdater) TestOverMaxRoyalty() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", nft.PaymentParameter(nft.MaxPaymentParameter+1), "https://localhost:5000/collection", []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "invalid range of paymentparameter")
}

func (t *testCollectionPolicyUpdater) TestEmptyUri() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "", []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionPolicyUpdater.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionPolicyUpdater.Fact())
	t.Implements((*operation.Operation)(nil), collectionPolicyUpdater)
}

func (t *testCollectionPolicyUpdater) TestSpaceUri() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "     ", []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "uri with only spaces")
}

func (t *testCollectionPolicyUpdater) TestNotTrimmedUri() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	uri := "   https://localhost:5000/collection   "
	policy := NewCollectionPolicy("New Collection", 0, nft.URI(uri), []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	_, compareErr := url.Parse(uri)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), compareErr.Error())
}

func (t *testCollectionPolicyUpdater) TestOverMaxUri() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	uri := strings.Repeat("a", nft.MaxURILength+1)
	policy := NewCollectionPolicy("New Collection", 0, nft.URI(uri), []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "invalid length of uri")
}

func (t *testCollectionPolicyUpdater) TestOverMaxWhites() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	policy := NewCollectionPolicy("New Collection", 0, "https://localhost:5000/collection", []base.Address{
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
		MustAddress(util.UUID().String()),
	})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "address in white list over allowed")
}

func (t *testCollectionPolicyUpdater) TestEmptyWhites() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "https://localhost:5000/collection", []base.Address{})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionPolicyUpdater.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionPolicyUpdater.Fact())
	t.Implements((*operation.Operation)(nil), collectionPolicyUpdater)
}

func (t *testCollectionPolicyUpdater) TestDuplicateWhites() {
	sender := MustAddress(util.UUID().String())
	white := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "https://localhost:5000/collection", []base.Address{white, white})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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

	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "duplicate white found")
}

func (t *testCollectionPolicyUpdater) TestOverSizeMemo() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	policy := NewCollectionPolicy("New Collection", 0, "https://localhost:5000/collection", []base.Address{sender})
	fact := NewCollectionPolicyUpdaterFact(token, sender, extensioncurrency.ContractID("ABC"), policy, "MCC")

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
	collectionPolicyUpdater, err := NewCollectionPolicyUpdater(fact, fs, memo)
	t.NoError(err)

	err = collectionPolicyUpdater.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestCollectionPolicyUpdaters(t *testing.T) {
	suite.Run(t, new(testCollectionPolicyUpdater))
}
