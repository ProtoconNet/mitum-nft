package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	MintFormType   = hint.Type("mitum-nft-mint-form")
	MintFormHint   = hint.NewHint(MintFormType, "v0.0.1")
	MintFormHinter = MintForm{BaseHinter: hint.NewBaseHinter(MintFormHint)}
)

type MintForm struct {
	hint.BaseHinter
	hash        nft.NFTHash
	uri         nft.NFTUri
	copyrighter nft.Copyrighter
}

func NewMintForm(hash nft.NFTHash, uri nft.NFTUri, copyrighter nft.Copyrighter) MintForm {
	return MintForm{
		BaseHinter:  hint.NewBaseHinter(MintFormHint),
		hash:        hash,
		uri:         uri,
		copyrighter: copyrighter,
	}
}

func MustNewMintform(hash nft.NFTHash, uri nft.NFTUri, copyrighter nft.Copyrighter) MintForm {
	form := NewMintForm(hash, uri, copyrighter)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form MintForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.hash.Bytes(),
		form.uri.Bytes(),
		form.copyrighter.Bytes(),
	)
}

func (form MintForm) IsValid([]byte) error {
	if err := form.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.hash,
		form.uri,
		form.copyrighter); err != nil {
		return err
	}

	return nil
}

var (
	MintFactType   = hint.Type("mitum-nft-mint-operation-fact")
	MintFactHint   = hint.NewHint(MintFactType, "v0.0.1")
	MintFactHinter = MintFact{BaseHinter: hint.NewBaseHinter(MintFactHint)}
	MintType       = hint.Type("mitum-nft-mint-operation")
	MintHint       = hint.NewHint(MintType, "v0.0.1")
	MintHinter     = Mint{BaseOperation: operationHinter(MintHint)}
)

type MintFact struct {
	hint.BaseHinter
	h          valuehash.Hash
	token      []byte
	sender     base.Address
	collection nft.Symbol
	form       MintForm
	cid        currency.CurrencyID
}

func NewMintFact(token []byte, sender base.Address, collection nft.Symbol, form MintForm, cid currency.CurrencyID) MintFact {
	fact := MintFact{
		BaseHinter: hint.NewBaseHinter(MintFactHint),
		token:      token,
		sender:     sender,
		collection: collection,
		form:       form,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact MintFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact MintFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.collection.Bytes(),
		fact.form.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact MintFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for MintFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.collection,
		fact.form,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact MintFact) Token() []byte {
	return fact.token
}

func (fact MintFact) Sender() base.Address {
	return fact.sender
}

func (fact MintFact) Collection() nft.Symbol {
	return fact.collection
}

func (fact MintFact) Form() MintForm {
	return fact.form
}

func (fact MintFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.Sender()

	return as, nil
}

func (fact MintFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact MintFact) Rebuild() MintFact {
	fact.h = fact.GenerateHash()

	return fact
}

type Mint struct {
	currency.BaseOperation
}

func NewMint(fact MintFact, fs []base.FactSign, memo string) (Mint, error) {
	bo, err := currency.NewBaseOperationFromFact(MintHint, fact, fs, memo)
	if err != nil {
		return Mint{}, err
	}
	return Mint{BaseOperation: bo}, nil
}
