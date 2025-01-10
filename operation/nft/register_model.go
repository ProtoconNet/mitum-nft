package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterModelFactHint = hint.MustNewHint("mitum-nft-register-model-operation-fact-v0.0.1")
	RegisterModelHint     = hint.MustNewHint("mitum-nft-register-model-operation-v0.0.1")
)

type RegisterModelFact struct {
	base.BaseFact
	sender          base.Address
	contract        base.Address
	name            types.CollectionName
	royalty         types.PaymentParameter
	uri             types.URI
	minterWhitelist []base.Address
	currency        ctypes.CurrencyID
}

func NewRegisterModelFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []base.Address,
	currency ctypes.CurrencyID,
) RegisterModelFact {
	bf := base.NewBaseFact(RegisterModelFactHint, token)
	fact := RegisterModelFact{
		BaseFact:        bf,
		sender:          sender,
		contract:        contract,
		name:            name,
		royalty:         royalty,
		uri:             uri,
		minterWhitelist: whitelist,
		currency:        currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RegisterModelFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.name,
		fact.royalty,
		fact.uri,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if l := len(fact.minterWhitelist); l > types.MaxWhitelist {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("whitelist over allowed, %d > %d", l, types.MaxWhitelist)))
	}

	founds := map[string]struct{}{}
	for _, white := range fact.minterWhitelist {
		if err := white.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if white.Equal(fact.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("whitelist %v is same with contract account", white)))
		}

		if _, found := founds[white.String()]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("whitelist %v", white)))
		}
		founds[white.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact RegisterModelFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RegisterModelFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterModelFact) Bytes() []byte {
	as := make([][]byte, len(fact.minterWhitelist))
	for i, white := range fact.minterWhitelist {
		as[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.name.Bytes(),
		fact.royalty.Bytes(),
		fact.uri.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (fact RegisterModelFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RegisterModelFact) Sender() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) Contract() base.Address {
	return fact.contract
}

func (fact RegisterModelFact) Name() types.CollectionName {
	return fact.name
}

func (fact RegisterModelFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact RegisterModelFact) URI() types.URI {
	return fact.uri
}

func (fact RegisterModelFact) WhiteList() []base.Address {
	return fact.minterWhitelist
}

func (fact RegisterModelFact) Addresses() ([]base.Address, error) {
	l := 2 + len(fact.minterWhitelist)

	as := make([]base.Address, l)
	copy(as, fact.minterWhitelist)

	as[l-2] = fact.sender
	as[l-1] = fact.contract

	return as, nil
}

func (fact RegisterModelFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact RegisterModelFact) FeePayer() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) FactUser() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) Signer() base.Address {
	return fact.sender
}

func (fact RegisterModelFact) InActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

func (fact RegisterModelFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

type RegisterModel struct {
	extras.ExtendedOperation
}

func NewRegisterModel(fact RegisterModelFact) (RegisterModel, error) {
	return RegisterModel{
		ExtendedOperation: extras.NewExtendedOperation(RegisterModelHint, fact),
	}, nil
}
