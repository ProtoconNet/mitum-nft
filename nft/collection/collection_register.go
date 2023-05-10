package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var CollectionRegisterFormHint = hint.MustNewHint("mitum-nft-collection-register-form-v0.0.1")

type CollectionRegisterForm struct {
	hint.BaseHinter
	contract   base.Address
	collection extensioncurrency.ContractID
	name       CollectionName
	royalty    nft.PaymentParameter
	uri        nft.URI
	whites     []base.Address
}

func NewCollectionRegisterForm(
	contract base.Address,
	collection extensioncurrency.ContractID,
	name CollectionName,
	royalty nft.PaymentParameter,
	uri nft.URI,
	whites []base.Address,
) CollectionRegisterForm {
	return CollectionRegisterForm{
		BaseHinter: hint.NewBaseHinter(CollectionRegisterFormHint),
		contract:   contract,
		collection: collection,
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whites:     whites,
	}
}

func (form CollectionRegisterForm) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		form.BaseHinter,
		form.contract,
		form.collection,
		form.name,
		form.royalty,
		form.uri,
	); err != nil {
		return err
	}

	if l := len(form.whites); l > MaxWhites {
		return util.ErrInvalid.Errorf("whites over allowed, %d > %d", l, MaxWhites)
	}

	founds := map[string]struct{}{}
	for _, white := range form.whites {
		if err := white.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[white.String()]; found {
			return util.ErrInvalid.Errorf("duplicate white found, %q", white)
		}
		founds[white.String()] = struct{}{}
	}

	return nil
}

func (form CollectionRegisterForm) Bytes() []byte {
	as := make([][]byte, len(form.whites))
	for i, white := range form.whites {
		as[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		form.contract.Bytes(),
		form.collection.Bytes(),
		form.name.Bytes(),
		form.royalty.Bytes(),
		form.uri.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (form CollectionRegisterForm) Contract() base.Address {
	return form.contract
}

func (form CollectionRegisterForm) Collection() extensioncurrency.ContractID {
	return form.collection
}

func (form CollectionRegisterForm) Name() CollectionName {
	return form.name
}

func (form CollectionRegisterForm) Royalty() nft.PaymentParameter {
	return form.royalty
}

func (form CollectionRegisterForm) URI() nft.URI {
	return form.uri
}

func (form CollectionRegisterForm) Whites() []base.Address {
	return form.whites
}

func (form CollectionRegisterForm) Addresses() ([]base.Address, error) {
	l := 1 + len(form.whites)

	as := make([]base.Address, l)
	copy(as, form.whites)

	as[l-1] = form.contract

	return as, nil
}

var (
	CollectionRegisterFactHint = hint.MustNewHint("mitum-nft-collection-register-operation-fact-v0.0.1")
	CollectionRegisterHint     = hint.MustNewHint("mitum-nft-collection-register-operation-v0.0.1")
)

type CollectionRegisterFact struct {
	base.BaseFact
	sender   base.Address
	form     CollectionRegisterForm
	currency currency.CurrencyID
}

func NewCollectionRegisterFact(token []byte, sender base.Address, form CollectionRegisterForm, currency currency.CurrencyID) CollectionRegisterFact {
	bf := base.NewBaseFact(CollectionRegisterFactHint, token)
	fact := CollectionRegisterFact{
		BaseFact: bf,
		sender:   sender,
		form:     form,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CollectionRegisterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.form,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.form.contract) {
		return util.ErrInvalid.Errorf("sender and contract are the same, %q == %q", fact.sender, fact.form.contract)
	}

	return nil
}

func (fact CollectionRegisterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CollectionRegisterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionRegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.form.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CollectionRegisterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CollectionRegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionRegisterFact) Form() CollectionRegisterForm {
	return fact.form
}

func (fact CollectionRegisterFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender}, nil
}

func (fact CollectionRegisterFact) Currency() currency.CurrencyID {
	return fact.currency
}

type CollectionRegister struct {
	currency.BaseOperation
}

func NewCollectionRegister(fact CollectionRegisterFact) (CollectionRegister, error) {
	return CollectionRegister{BaseOperation: currency.NewBaseOperation(CollectionRegisterHint, fact)}, nil
}

func (op *CollectionRegister) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
