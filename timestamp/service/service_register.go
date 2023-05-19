package service

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	ServiceRegisterFactHint = hint.MustNewHint("mitum-timestamp-service-register-operation-fact-v0.0.1")
	ServiceRegisterHint     = hint.MustNewHint("mitum-timestamp-service-register-operation-v0.0.1")
)

type ServiceRegisterFact struct {
	base.BaseFact
	sender   base.Address
	target   base.Address
	service  extensioncurrency.ContractID
	currency currency.CurrencyID
}

func NewServiceRegisterFact(token []byte, sender, target base.Address, service extensioncurrency.ContractID, currency currency.CurrencyID) ServiceRegisterFact {
	bf := base.NewBaseFact(ServiceRegisterFactHint, token)
	fact := ServiceRegisterFact{
		BaseFact: bf,
		sender:   sender,
		target:   target,
		service:  service,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ServiceRegisterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.target,
		fact.service,
		fact.currency,
	); err != nil {
		return err
	}

	return nil
}

func (fact ServiceRegisterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ServiceRegisterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ServiceRegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.service.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact ServiceRegisterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ServiceRegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact ServiceRegisterFact) Target() base.Address {
	return fact.target
}

func (fact ServiceRegisterFact) Service() extensioncurrency.ContractID {
	return fact.service
}

func (fact ServiceRegisterFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender, fact.target}, nil
}

func (fact ServiceRegisterFact) Currency() currency.CurrencyID {
	return fact.currency
}

type ServiceRegister struct {
	currency.BaseOperation
}

func NewServiceRegister(fact ServiceRegisterFact) (ServiceRegister, error) {
	return ServiceRegister{BaseOperation: currency.NewBaseOperation(ServiceRegisterHint, fact)}, nil
}

func (op *ServiceRegister) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
