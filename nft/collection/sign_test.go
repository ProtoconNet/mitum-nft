package collection

import (
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

type testSign struct {
	suite.Suite
}

func (t *testSign) TestNew() {
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

func (t *testSign) TestDeplicateNFTID() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)

	items := []SignItem{
		NewSignItem(CreatorQualification, nid, "MCC"),
		NewSignItem(CopyrighterQualification, nid, "MCC"),
	}
	fact := NewSignFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "duplicate nft found")
}

func (t *testSign) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []SignItem{}
	fact := NewSignFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "empty items for SignFact")
}

func (t *testSign) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []SignItem{
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1), "MCC"),
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 2), "MCC"),
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 3), "MCC"),
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 4), "MCC"),
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 5), "MCC"),
		NewSignItem(CreatorQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 6), "MCC"),
		NewSignItem(CopyrighterQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 7), "MCC"),
		NewSignItem(CopyrighterQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 8), "MCC"),
		NewSignItem(CopyrighterQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 9), "MCC"),
		NewSignItem(CopyrighterQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 10), "MCC"),
		NewSignItem(CopyrighterQualification, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 11), "MCC"),
	}
	fact := NewSignFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	sign, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testSign) TestOverSizeMemo() {
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

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	sign, err := NewSign(fact, fs, memo)
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestSign(t *testing.T) {
	suite.Run(t, new(testSign))
}
