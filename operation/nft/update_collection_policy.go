package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	UpdateCollectionPolicyFactHint = hint.MustNewHint("mitum-nft-update-collection-policy-operation-fact-v0.0.1")
	UpdateCollectionPolicyHint     = hint.MustNewHint("mitum-nft-update-collection-policy-operation-v0.0.1")
)

type UpdateCollectionPolicyFact struct {
	mitumbase.BaseFact
	sender    mitumbase.Address
	contract  mitumbase.Address
	name      types.CollectionName
	royalty   types.PaymentParameter
	uri       types.URI
	whitelist []mitumbase.Address
	currency  currencytypes.CurrencyID
}

func NewUpdateCollectionPolicyFact(
	token []byte,
	sender, contract mitumbase.Address,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []mitumbase.Address,
	currency currencytypes.CurrencyID,
) UpdateCollectionPolicyFact {
	bf := mitumbase.NewBaseFact(UpdateCollectionPolicyFactHint, token)

	fact := UpdateCollectionPolicyFact{
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

func (fact UpdateCollectionPolicyFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if l := len(fact.whitelist); l > types.MaxWhitelist {
		return util.ErrInvalid.Errorf("whitelist over allowed, %d > %d", l, types.MaxWhitelist)
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
		return err
	}

	founds := map[string]struct{}{}
	for _, white := range fact.whitelist {
		if err := white.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[white.String()]; found {
			return util.ErrInvalid.Errorf("duplicate whitelist account found, %q", white)
		}
		founds[white.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact UpdateCollectionPolicyFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdateCollectionPolicyFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateCollectionPolicyFact) Bytes() []byte {
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

func (fact UpdateCollectionPolicyFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact UpdateCollectionPolicyFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact UpdateCollectionPolicyFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact UpdateCollectionPolicyFact) Name() types.CollectionName {
	return fact.name
}

func (fact UpdateCollectionPolicyFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact UpdateCollectionPolicyFact) URI() types.URI {
	return fact.uri
}

func (fact UpdateCollectionPolicyFact) Whitelist() []mitumbase.Address {
	return fact.whitelist
}

func (fact UpdateCollectionPolicyFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact UpdateCollectionPolicyFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type UpdateCollectionPolicy struct {
	common.BaseOperation
}

func NewUpdateCollectionPolicy(fact UpdateCollectionPolicyFact) (UpdateCollectionPolicy, error) {
	return UpdateCollectionPolicy{BaseOperation: common.NewBaseOperation(UpdateCollectionPolicyHint, fact)}, nil
}
