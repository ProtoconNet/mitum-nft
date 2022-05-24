package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

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
	target base.Address
	design nft.Design
	cid    currency.CurrencyID
}

func NewCollectionRegisterFact(token []byte, sender base.Address, target base.Address, design nft.Design, cid currency.CurrencyID) CollectionRegisterFact {
	fact := CollectionRegisterFact{
		BaseHinter: hint.NewBaseHinter(CollectionRegisterFactHint),
		token:      token,
		sender:     sender,
		target:     target,
		design:     design,
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
		fact.target.Bytes(),
		fact.design.Bytes(),
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
		fact.target,
		fact.design,
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

func (fact CollectionRegisterFact) Target() base.Address {
	return fact.target
}

func (fact CollectionRegisterFact) Design() nft.Design {
	return fact.design
}

func (fact CollectionRegisterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.target

	return as, nil
}

func (fact CollectionRegisterFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact CollectionRegisterFact) Rebuild() CollectionRegisterFact {
	design := fact.design.Rebuild()
	fact.design = design

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
