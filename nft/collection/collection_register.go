package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	CollectionRegisterFormType   = hint.Type("mitum-nft-collection-register-form")
	CollectionRegisterFormHint   = hint.NewHint(CollectionRegisterFormType, "v0.0.1")
	CollectionRegisterFormHinter = CollectionRegisterForm{BaseHinter: hint.NewBaseHinter(CollectionRegisterFormHint)}
)

type CollectionRegisterForm struct {
	hint.BaseHinter
	target  base.Address
	symbol  extensioncurrency.ContractID
	name    CollectionName
	royalty nft.PaymentParameter
	uri     nft.URI
	whites  []base.Address
}

func NewCollectionRegisterForm(target base.Address, symbol extensioncurrency.ContractID, name CollectionName,
	royalty nft.PaymentParameter, uri nft.URI, whites []base.Address) CollectionRegisterForm {
	return CollectionRegisterForm{
		BaseHinter: hint.NewBaseHinter(CollectionRegisterFormHint),
		target:     target,
		symbol:     symbol,
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whites:     whites,
	}
}

func MustNewCollectionRegisterForm(target base.Address, symbol extensioncurrency.ContractID, name CollectionName,
	royalty nft.PaymentParameter, uri nft.URI, whites []base.Address) CollectionRegisterForm {
	form := NewCollectionRegisterForm(target, symbol, name, royalty, uri, whites)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form CollectionRegisterForm) Bytes() []byte {
	as := make([][]byte, len(form.whites))
	for i := range form.whites {
		as[i] = form.whites[i].Bytes()
	}

	return util.ConcatBytesSlice(
		form.target.Bytes(),
		form.symbol.Bytes(),
		form.name.Bytes(),
		form.royalty.Bytes(),
		form.uri.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (form CollectionRegisterForm) Target() base.Address {
	return form.target
}

func (form CollectionRegisterForm) Symbol() extensioncurrency.ContractID {
	return form.symbol
}

func (form CollectionRegisterForm) Name() CollectionName {
	return form.name
}

func (form CollectionRegisterForm) Royalty() nft.PaymentParameter {
	return form.royalty
}

func (form CollectionRegisterForm) Uri() nft.URI {
	return form.uri
}

func (form CollectionRegisterForm) Whites() []base.Address {
	return form.whites
}

func (form CollectionRegisterForm) Addresses() ([]base.Address, error) {
	l := 1 + len(form.whites)

	as := make([]base.Address, l)
	for i := range form.whites {
		as[i] = form.whites[i]
	}

	as[l-1] = form.target

	return as, nil
}

func (form CollectionRegisterForm) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		form.BaseHinter,
		form.target,
		form.symbol,
		form.name,
		form.royalty,
		form.uri,
	); err != nil {
		return err
	}

	if l := len(form.whites); l > MaxWhiteAddress {
		return isvalid.InvalidError.Errorf("address in white list over allowed; %d > %d", l, MaxWhiteAddress)
	}

	founds := map[base.Address]struct{}{}
	for i := range form.whites {
		acc := form.whites[i]
		if err := acc.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[acc]; found {
			return isvalid.InvalidError.Errorf("duplicate white found; %q", acc)
		}
		founds[acc] = struct{}{}
	}

	return nil
}

func (form CollectionRegisterForm) Rebuild() CollectionRegisterForm {
	return form
}

var (
	CollectionRegisterFactType   = hint.Type("mitum-nft-collection-register-operation-fact")
	CollectionRegisterFactHint   = hint.NewHint(CollectionRegisterFactType, "v0.0.1")
	CollectionRegisterFactHinter = CollectionRegisterFact{BaseHinter: hint.NewBaseHinter(CollectionRegisterFactHint)}
	CollectionRegisterType       = hint.Type("mitum-nft-collection-register-operation")
	CollectionRegisterHint       = hint.NewHint(CollectionRegisterType, "v0.0.1")
	CollectionRegisterHinter     = CollectionRegister{BaseOperation: operationHinter(CollectionRegisterHint)}
)

type CollectionRegisterFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	form   CollectionRegisterForm
	cid    currency.CurrencyID
}

func NewCollectionRegisterFact(token []byte, sender base.Address, form CollectionRegisterForm, cid currency.CurrencyID) CollectionRegisterFact {
	fact := CollectionRegisterFact{
		BaseHinter: hint.NewBaseHinter(CollectionRegisterFactHint),
		token:      token,
		sender:     sender,
		form:       form,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact CollectionRegisterFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact CollectionRegisterFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionRegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.form.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact CollectionRegisterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for CollectionRegisterFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.form,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact CollectionRegisterFact) Token() []byte {
	return fact.token
}

func (fact CollectionRegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionRegisterFact) Form() CollectionRegisterForm {
	return fact.form
}

func (fact CollectionRegisterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact CollectionRegisterFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact CollectionRegisterFact) Rebuild() CollectionRegisterFact {
	form := fact.form.Rebuild()
	fact.form = form

	fact.h = fact.GenerateHash()

	return fact
}

type CollectionRegister struct {
	currency.BaseOperation
}

func NewCollectionRegister(fact CollectionRegisterFact, fs []base.FactSign, memo string) (CollectionRegister, error) {
	bo, err := currency.NewBaseOperationFromFact(CollectionRegisterHint, fact, fs, memo)
	if err != nil {
		return CollectionRegister{}, err
	}
	return CollectionRegister{BaseOperation: bo}, nil
}
