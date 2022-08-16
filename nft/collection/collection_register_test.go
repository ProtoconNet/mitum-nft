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

type testCollectionRegister struct {
	suite.Suite
}

func (t *testCollectionRegister) TestNew() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "https://localhost:5000/collection", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionRegister.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionRegister.Fact())
	t.Implements((*operation.Operation)(nil), collectionRegister)
}

func (t *testCollectionRegister) TestSameSenderTarget() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		sender, extensioncurrency.ContractID("ABC"), "Collection", 0, "https://localhost:5000/collection", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "sender and target are the same")
}

func (t *testCollectionRegister) TestOverMaxRoyalty() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", nft.PaymentParameter(nft.MaxPaymentParameter+1), "https://localhost:5000/collection", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "invalid range of paymentparameter")
}

func (t *testCollectionRegister) TestEmptyUri() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionRegister.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionRegister.Fact())
	t.Implements((*operation.Operation)(nil), collectionRegister)
}

func (t *testCollectionRegister) TestSpaceUri() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "      ", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "uri with only spaces")
}

func (t *testCollectionRegister) TestNotTrimmedUri() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	uri := "   https://localhost:5000/collection   "
	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, nft.URI(uri), []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	_, compareErr := url.Parse(uri)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), compareErr.Error())
}

func (t *testCollectionRegister) TestOverMaxUri() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	uri := strings.Repeat("a", nft.MaxURILength+1)
	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, nft.URI(uri), []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "invalid length of uri")
}

func (t *testCollectionRegister) TestOverMaxWhites() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "https://localhost:5000/collection", []base.Address{
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
		},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "address in white list over allowed")
}

func (t *testCollectionRegister) TestEmptyWhites() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "https://localhost:5000/collection", []base.Address{},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	t.NoError(collectionRegister.IsValid(nil))

	t.Implements((*base.Fact)(nil), collectionRegister.Fact())
	t.Implements((*operation.Operation)(nil), collectionRegister)
}

func (t *testCollectionRegister) TestDuplicateWhites() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())
	white := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"),
		"Collection",
		0,
		"https://localhost:5000/collection",
		[]base.Address{white, white},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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

	collectionRegister, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "duplicate white found")
}

func (t *testCollectionRegister) TestOverSizeMemo() {
	sender := MustAddress(util.UUID().String())
	target := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	form := NewCollectionRegisterForm(
		target, extensioncurrency.ContractID("ABC"), "Collection", 0, "https://localhost:5000/collection", []base.Address{sender},
	)
	fact := NewCollectionRegisterFact(token, sender, form, "MCC")

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
	collectionRegister, err := NewCollectionRegister(fact, fs, memo)
	t.NoError(err)

	err = collectionRegister.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestCollectionRegisters(t *testing.T) {
	suite.Run(t, new(testCollectionRegister))
}
