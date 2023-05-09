package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CollectionPolicyUpdaterFactHint = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-fact-v0.0.1")
	CollectionPolicyUpdaterHint     = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-v0.0.1")
)

type CollectionPolicyUpdaterFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	collection extensioncurrency.ContractID
	policy     CollectionPolicy
	currency   currency.CurrencyID
}

func NewCollectionPolicyUpdaterFact(
	token []byte,
	sender, contract base.Address,
	collection extensioncurrency.ContractID,
	policy CollectionPolicy,
	currency currency.CurrencyID,
) CollectionPolicyUpdaterFact {
	bf := base.NewBaseFact(CollectionPolicyUpdaterFactHint, token)

	fact := CollectionPolicyUpdaterFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		collection: collection,
		policy:     policy,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CollectionPolicyUpdaterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(
		nil, false,
		fact.sender,
		fact.contract,
		fact.collection,
		fact.policy,
		fact.currency,
	); err != nil {
		return err
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
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.collection.Bytes(),
		fact.policy.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CollectionPolicyUpdaterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CollectionPolicyUpdaterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionPolicyUpdaterFact) Contract() base.Address {
	return fact.contract
}

func (fact CollectionPolicyUpdaterFact) Collection() extensioncurrency.ContractID {
	return fact.collection
}

func (fact CollectionPolicyUpdaterFact) Policy() CollectionPolicy {
	return fact.policy
}

func (fact CollectionPolicyUpdaterFact) Currency() currency.CurrencyID {
	return fact.currency
}

func (fact CollectionPolicyUpdaterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

type CollectionPolicyUpdater struct {
	currency.BaseOperation
}

func NewCollectionPolicyUpdater(fact CollectionPolicyUpdaterFact) (CollectionPolicyUpdater, error) {
	return CollectionPolicyUpdater{BaseOperation: currency.NewBaseOperation(CollectionPolicyUpdaterHint, fact)}, nil
}

func (op *CollectionPolicyUpdater) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
