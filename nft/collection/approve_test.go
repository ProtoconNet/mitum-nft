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

type testApprove struct {
	suite.Suite
}

func (t *testApprove) TestNew() {
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

func (t *testApprove) TestDeplicateNFTID() {
	sender := MustAddress(util.UUID().String())
	approved0 := MustAddress(util.UUID().String())
	approved1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []ApproveItem{
		NewApproveItem(approved0, nid, "MCC"),
		NewApproveItem(approved1, nid, "MCC"),
	}
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
	t.Contains(err.Error(), "duplicate nft found")
}

func (t *testApprove) TestSameSenderApproved() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []ApproveItem{NewApproveItem(sender, nid, "MCC")}
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

func (t *testApprove) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []ApproveItem{}
	fact := NewApproveFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = approve.IsValid(nil)
	t.Contains(err.Error(), "empty items for ApproveFact")
}

func (t *testApprove) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []ApproveItem{
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(1), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(2), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(3), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(4), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(5), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(6), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(7), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(8), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(9), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(10), "MCC"),
		NewApproveItem(MustAddress(util.UUID().String()), nft.NewTestNFTID(11), "MCC"),
	}
	fact := NewApproveFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = approve.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testApprove) TestOverSizeMemo() {
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

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	approve, err := NewApprove(fact, fs, memo)
	t.NoError(err)

	err = approve.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestApproves(t *testing.T) {
	suite.Run(t, new(testApprove))
}
