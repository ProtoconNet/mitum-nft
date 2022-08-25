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

type testMint struct {
	suite.Suite
}

func (t *testMint) TestNew() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 50, false),
			nft.NewSigner(creator1, 50, false),
		},
	)
	copyrighters := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(copyrighter0, 50, false),
			nft.NewSigner(copyrighter1, 50, false),
		},
	)

	form := NewMintForm(
		nft.NFTHash(nft.NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft", creators, copyrighters,
	)
	items := []MintItem{
		NewMintItem(extensioncurrency.ContractID("ABC"), form, "MCC"),
	}
	fact := NewMintFact(token, sender, items)

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

	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	t.NoError(mint.IsValid(nil))

	t.Implements((*base.Fact)(nil), mint.Fact())
	t.Implements((*operation.Operation)(nil), mint)
}

func (t *testMint) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []MintItem{}
	fact := NewMintFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "empty items for MintFact")
}

func (t *testMint) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 50, false),
			nft.NewSigner(creator1, 50, false),
		},
	)
	copyrighters := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(copyrighter0, 50, false),
			nft.NewSigner(copyrighter1, 50, false),
		},
	)

	forms := []MintForm{
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(1).Hash().String()), "https://localhost:5000/nft/1", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(2).Hash().String()), "https://localhost:5000/nft/2", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(3).Hash().String()), "https://localhost:5000/nft/3", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(4).Hash().String()), "https://localhost:5000/nft/4", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(5).Hash().String()), "https://localhost:5000/nft/5", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(6).Hash().String()), "https://localhost:5000/nft/6", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(7).Hash().String()), "https://localhost:5000/nft/7", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(8).Hash().String()), "https://localhost:5000/nft/8", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(9).Hash().String()), "https://localhost:5000/nft/9", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(10).Hash().String()), "https://localhost:5000/nft/10", creators, copyrighters),
		NewMintForm(nft.NFTHash(nft.NewTestNFTID(11).Hash().String()), "https://localhost:5000/nft/11", creators, copyrighters),
	}

	items := []MintItem{
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[0], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[1], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[2], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[3], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[4], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[5], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[6], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[7], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[8], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[9], "MCC"),
		NewMintItem(extensioncurrency.ContractID("ABC"), forms[10], "MCC"),
	}
	fact := NewMintFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testMint) TestOverSizeMemo() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 50, false),
			nft.NewSigner(creator1, 50, false),
		},
	)
	copyrighters := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(copyrighter0, 50, false),
			nft.NewSigner(copyrighter1, 50, false),
		},
	)

	form := NewMintForm(
		nft.NFTHash(nft.NewTestNFTID(1).Hash().String()),
		"https://localhost:5000/nft", creators, copyrighters,
	)
	items := []MintItem{
		NewMintItem(extensioncurrency.ContractID("ABC"), form, "MCC"),
	}
	fact := NewMintFact(token, sender, items)

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
	mint, err := NewMint(fact, fs, memo)
	t.NoError(err)

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestMint(t *testing.T) {
	suite.Run(t, new(testMint))
}
