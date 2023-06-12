package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CollectionPolicyUpdaterFactHint = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-fact-v0.0.1")
	CollectionPolicyUpdaterHint     = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-v0.0.1")
)

type CollectionPolicyUpdaterFact struct {
	mitumbase.BaseFact
	sender     mitumbase.Address
	contract   mitumbase.Address
	collection currencytypes.ContractID
	name       types.CollectionName
	royalty    types.PaymentParameter
	uri        types.URI
	whitelist  []mitumbase.Address
	currency   currencytypes.CurrencyID
}

func NewCollectionPolicyUpdaterFact(
	token []byte,
	sender, contract mitumbase.Address,
	collection currencytypes.ContractID,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []mitumbase.Address,
	currency currencytypes.CurrencyID,
) CollectionPolicyUpdaterFact {
	bf := mitumbase.NewBaseFact(CollectionPolicyUpdaterFactHint, token)

	fact := CollectionPolicyUpdaterFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		collection: collection,
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whitelist:  whitelist,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CollectionPolicyUpdaterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.whitelist); l > types.MaxWhites {
		return util.ErrInvalid.Errorf("whites over allowed, %d > %d", l, types.MaxWhites)
	}

	if err := util.CheckIsValiders(
		nil, false,
		fact.sender,
		fact.contract,
		fact.collection,
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

	return nil
}

func (fact CollectionPolicyUpdaterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CollectionPolicyUpdaterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionPolicyUpdaterFact) Bytes() []byte {
	as := make([][]byte, len(fact.whitelist))
	for i, white := range fact.whitelist {
		as[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.collection.Bytes(),
		fact.name.Bytes(),
		fact.royalty.Bytes(),
		fact.uri.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (fact CollectionPolicyUpdaterFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CollectionPolicyUpdaterFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CollectionPolicyUpdaterFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact CollectionPolicyUpdaterFact) Collection() currencytypes.ContractID {
	return fact.collection
}

func (fact CollectionPolicyUpdaterFact) Name() types.CollectionName {
	return fact.name
}

func (fact CollectionPolicyUpdaterFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact CollectionPolicyUpdaterFact) URI() types.URI {
	return fact.uri
}

func (fact CollectionPolicyUpdaterFact) Whitelist() []mitumbase.Address {
	return fact.whitelist
}

func (fact CollectionPolicyUpdaterFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CollectionPolicyUpdaterFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type CollectionPolicyUpdater struct {
	common.BaseOperation
}

func NewCollectionPolicyUpdater(fact CollectionPolicyUpdaterFact) (CollectionPolicyUpdater, error) {
	return CollectionPolicyUpdater{BaseOperation: common.NewBaseOperation(CollectionPolicyUpdaterHint, fact)}, nil
}

func (op *CollectionPolicyUpdater) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
