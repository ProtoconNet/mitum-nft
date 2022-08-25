package collection

import (
	"net/url"
	"strings"
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

type testMintForm struct {
	suite.Suite
}

func (t *testMintForm) TestNew() {
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

func (t *testMintForm) TestEmptyNFTHash() {
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
		"",
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

func (t *testMintForm) TestOverMaxNFTHash() {
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

	hash := strings.Repeat("a", nft.MaxNFTHashLength) + "a"
	form := NewMintForm(
		nft.NFTHash(hash),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "invalid length of nft hash")
}

func (t *testMintForm) TestSpaceNFTHash() {
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
		"           ",
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "nft hash with only spaces")
}

func (t *testMintForm) TestNotTrimmedNFTHash() {
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
		"   "+nft.NFTHash(nft.NewTestNFTID(1).Hash().String())+"   ",
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

func (t *testMintForm) TestEmptyUri() {
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
		"", creators, copyrighters,
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "empty uri")
}

func (t *testMintForm) TestOverMaxUri() {
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

	uri := strings.Repeat("a", nft.MaxURILength) + "a"
	form := NewMintForm(
		nft.NFTHash(nft.NewTestNFTID(1).Hash().String()),
		nft.URI(uri), creators, copyrighters,
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "invalid length of uri")
}

func (t *testMintForm) TestSpaceUri() {
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
		"           ", creators, copyrighters,
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "uri with only spaces")
}

func (t *testMintForm) TestNotTrimmedUri() {
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

	uri := "   https://localhost:5000/nft   "
	form := NewMintForm(
		nft.NFTHash(nft.NewTestNFTID(1).Hash().String()),
		nft.URI(uri), creators, copyrighters,
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

	_, compareErr := url.Parse(uri)

	err = mint.IsValid(nil)
	t.Contains(err.Error(), compareErr.Error())
}

func (t *testMintForm) TestEmptyCreators() {
	sender := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		0, []nft.Signer{},
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

func (t *testMintForm) TestOverMaxCreators() {
	sender := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(MustAddress(util.UUID().String()), 0, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "signers over allowed")
}

func (t *testMintForm) TestEmptyCopyrighters() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	copyrighters := nft.NewSigners(
		0, []nft.Signer{},
	)
	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 50, false),
			nft.NewSigner(creator1, 50, false),
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

func (t *testMintForm) TestOverMaxCopyrighters() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	copyrighters := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(MustAddress(util.UUID().String()), 0, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
			nft.NewSigner(MustAddress(util.UUID().String()), 10, false),
		},
	)
	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 50, false),
			nft.NewSigner(creator1, 50, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "signers over allowed")
}

func (t *testMintForm) TestWrongCreatorsTotal() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		90, []nft.Signer{
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "total share must be equal to the sum of all shares")
}

func (t *testMintForm) TestWrongCopyrightersTotal() {
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
		90, []nft.Signer{
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "total share must be equal to the sum of all shares")
}

func (t *testMintForm) TestZeroCreatorsTotal() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		0, []nft.Signer{
			nft.NewSigner(creator0, 0, false),
			nft.NewSigner(creator1, 0, false),
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

func (t *testMintForm) TestZeroCopyrightersTotal() {
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
		0, []nft.Signer{
			nft.NewSigner(copyrighter0, 0, false),
			nft.NewSigner(copyrighter1, 0, false),
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

func (t *testMintForm) TestOverMaxCreatorsTotal() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		120, []nft.Signer{
			nft.NewSigner(creator0, 60, false),
			nft.NewSigner(creator1, 60, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "total share is over max")
}

func (t *testMintForm) TestOverMaxCopyrightersTotal() {
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
		120, []nft.Signer{
			nft.NewSigner(copyrighter0, 60, false),
			nft.NewSigner(copyrighter1, 60, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "total share is over max")
}

func (t *testMintForm) TestDuplicateCreator() {
	sender := MustAddress(util.UUID().String())
	creator0 := MustAddress(util.UUID().String())
	creator1 := MustAddress(util.UUID().String())
	copyrighter0 := MustAddress(util.UUID().String())
	copyrighter1 := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	creators := nft.NewSigners(
		100, []nft.Signer{
			nft.NewSigner(creator0, 30, false),
			nft.NewSigner(creator0, 30, false),
			nft.NewSigner(creator1, 40, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "duplicate signer found")
}

func (t *testMintForm) TestDuplicateCopyrighter() {
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
			nft.NewSigner(copyrighter0, 30, false),
			nft.NewSigner(copyrighter1, 30, false),
			nft.NewSigner(copyrighter1, 40, false),
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

	err = mint.IsValid(nil)
	t.Contains(err.Error(), "duplicate signer found")
}

func testMintFormEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
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

		return mint
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Mint)
		tb := b.(Mint)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(MintFact)
		ufact := tb.Fact().(MintFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.Equal(a.Collection(), b.Collection())
			t.Equal(a.Currency(), b.Currency())

			af := a.Form()
			bf := b.Form()
			t.Equal(af.NftHash(), bf.NftHash())
			t.Equal(af.Uri(), bf.Uri())
			t.True(af.Creators().Equal(bf.Creators()))
			t.True(af.Copyrighters().Equal(bf.Copyrighters()))
		}
	}

	return t
}

func TestMintFormEncodeJSON(t *testing.T) {
	suite.Run(t, testMintFormEncode(jsonenc.NewEncoder()))
}

func TestMintFormEncodeBSON(t *testing.T) {
	suite.Run(t, testMintFormEncode(bsonenc.NewEncoder()))
}

func TestMintForm(t *testing.T) {
	suite.Run(t, new(testMintForm))
}

type testMintItem struct {
	suite.Suite
}

// test creator qualification
func (t *testMintItem) TestNew() {
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

func TestMintItem(t *testing.T) {
	suite.Run(t, new(testSignItem))
}

func testMintItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
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

		return mint
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Mint)
		tb := b.(Mint)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(MintFact)
		ufact := tb.Fact().(MintFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.Equal(a.Collection(), b.Collection())
			t.Equal(a.Currency(), b.Currency())

			af := a.Form()
			bf := b.Form()

			t.Equal(af.NftHash(), bf.NftHash())
			t.Equal(af.Uri(), bf.Uri())
			t.True(af.Creators().Equal(bf.Creators()))
			t.True(af.Copyrighters().Equal(bf.Copyrighters()))
		}
	}

	return t
}

func TestMintItemEncodeJSON(t *testing.T) {
	suite.Run(t, testMintItemEncode(jsonenc.NewEncoder()))
}

func TestMintItemEncodeBSON(t *testing.T) {
	suite.Run(t, testMintItemEncode(bsonenc.NewEncoder()))
}
