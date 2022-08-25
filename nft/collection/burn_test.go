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

type testBurn struct {
	suite.Suite
}

func (t *testBurn) TestNew() {
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

func (t *testBurn) TestDeplicateNFTID() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)

	items := []BurnItem{
		NewBurnItem(nid, "MCC"),
		NewBurnItem(nid, "MCC"),
	}
	fact := NewBurnFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	err = burn.IsValid(nil)
	t.Contains(err.Error(), "duplicate nft found")
}

func (t *testBurn) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []BurnItem{}
	fact := NewBurnFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	err = burn.IsValid(nil)
	t.Contains(err.Error(), "empty items for BurnFact")
}

func (t *testBurn) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []BurnItem{
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 2), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 3), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 4), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 5), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 6), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 7), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 8), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 9), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 10), "MCC"),
		NewBurnItem(nft.NewNFTID(extensioncurrency.ContractID("ABC"), 11), "MCC"),
	}
	fact := NewBurnFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	burn, err := NewBurn(fact, fs, "")
	t.NoError(err)

	err = burn.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testBurn) TestOverSizeMemo() {
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

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	sign, err := NewBurn(fact, fs, memo)
	t.NoError(err)

	err = sign.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestBurn(t *testing.T) {
	suite.Run(t, new(testBurn))
}
