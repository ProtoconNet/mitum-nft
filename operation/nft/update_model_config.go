package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	UpdateModelConfigFactHint = hint.MustNewHint("mitum-nft-update-model-config-operation-fact-v0.0.1")
	UpdateModelConfigHint     = hint.MustNewHint("mitum-nft-update-model-config-operation-v0.0.1")
)

type UpdateModelConfigFact struct {
	mitumbase.BaseFact
	sender    mitumbase.Address
	contract  mitumbase.Address
	name      types.CollectionName
	royalty   types.PaymentParameter
	uri       types.URI
	whitelist []mitumbase.Address
	currency  currencytypes.CurrencyID
}

func NewUpdateModelConfigFact(
	token []byte,
	sender, contract mitumbase.Address,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []mitumbase.Address,
	currency currencytypes.CurrencyID,
) UpdateModelConfigFact {
	bf := mitumbase.NewBaseFact(UpdateModelConfigFactHint, token)

	fact := UpdateModelConfigFact{
		BaseFact:  bf,
		sender:    sender,
		contract:  contract,
		name:      name,
		royalty:   royalty,
		uri:       uri,
		whitelist: whitelist,
		currency:  currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact UpdateModelConfigFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract", fact.sender)))
	}

	if l := len(fact.whitelist); l > types.MaxWhitelist {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("whitelist over allowed, %d > %d", l, types.MaxWhitelist)))
	}

	if err := util.CheckIsValiders(
		nil, false,
		fact.sender,
		fact.contract,
		fact.name,
		fact.royalty,
		fact.uri,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, white := range fact.whitelist {
		if err := white.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if white.Equal(fact.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("whitelist account is same with contract")))
		}

		if _, found := founds[white.String()]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("whitelist account, %v", white)))
		}
		founds[white.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact UpdateModelConfigFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdateModelConfigFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateModelConfigFact) Bytes() []byte {
	as := make([][]byte, len(fact.whitelist))
	for i, white := range fact.whitelist {
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

func (fact UpdateModelConfigFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact UpdateModelConfigFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact UpdateModelConfigFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact UpdateModelConfigFact) Name() types.CollectionName {
	return fact.name
}

func (fact UpdateModelConfigFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact UpdateModelConfigFact) URI() types.URI {
	return fact.uri
}

func (fact UpdateModelConfigFact) Whitelist() []mitumbase.Address {
	return fact.whitelist
}

func (fact UpdateModelConfigFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact UpdateModelConfigFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type UpdateModelConfig struct {
	common.BaseOperation
}

func NewUpdateModelConfig(fact UpdateModelConfigFact) (UpdateModelConfig, error) {
	return UpdateModelConfig{BaseOperation: common.NewBaseOperation(UpdateModelConfigHint, fact)}, nil
}
