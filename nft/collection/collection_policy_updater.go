package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	CollectionPolicyUpdaterFactType   = hint.Type("mitum-nft-collection-policy-updater-operation-fact")
	CollectionPolicyUpdaterFactHint   = hint.NewHint(CollectionPolicyUpdaterFactType, "v0.0.1")
	CollectionPolicyUpdaterFactHinter = CollectionPolicyUpdaterFact{BaseHinter: hint.NewBaseHinter(CollectionPolicyUpdaterFactHint)}
	CollectionPolicyUpdaterType       = hint.Type("mitum-nft-collection-policy-updater-operation")
	CollectionPolicyUpdaterHint       = hint.NewHint(CollectionPolicyUpdaterType, "v0.0.1")
	CollectionPolicyUpdaterHinter     = CollectionPolicyUpdater{BaseOperation: operationHinter(CollectionPolicyUpdaterHint)}
)

type CollectionPolicyUpdaterFact struct {
	hint.BaseHinter
	h          valuehash.Hash
	token      []byte
	sender     base.Address
	collection extensioncurrency.ContractID
	policy     CollectionPolicy
	cid        currency.CurrencyID
}

func NewCollectionPolicyUpdaterFact(token []byte, sender base.Address, collection extensioncurrency.ContractID, policy CollectionPolicy, cid currency.CurrencyID) CollectionPolicyUpdaterFact {
	fact := CollectionPolicyUpdaterFact{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyUpdaterFactHint),
		token:      token,
		sender:     sender,
		collection: collection,
		policy:     policy,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact CollectionPolicyUpdaterFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact CollectionPolicyUpdaterFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionPolicyUpdaterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.collection.Bytes(),
		fact.policy.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact CollectionPolicyUpdaterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for CollectionPolicyUpdaterFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.collection,
		fact.policy,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact CollectionPolicyUpdaterFact) Token() []byte {
	return fact.token
}

func (fact CollectionPolicyUpdaterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionPolicyUpdaterFact) Collection() extensioncurrency.ContractID {
	return fact.collection
}

func (fact CollectionPolicyUpdaterFact) Policy() CollectionPolicy {
	return fact.policy
}

func (fact CollectionPolicyUpdaterFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact CollectionPolicyUpdaterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact CollectionPolicyUpdaterFact) Rebuild() CollectionPolicyUpdaterFact {
	policy := fact.policy.Rebuild().(CollectionPolicy)
	fact.policy = policy

	fact.h = fact.GenerateHash()

	return fact
}

type CollectionPolicyUpdater struct {
	currency.BaseOperation
}

func NewCollectionPolicyUpdater(fact CollectionPolicyUpdaterFact, fs []base.FactSign, memo string) (CollectionPolicyUpdater, error) {
	bo, err := currency.NewBaseOperationFromFact(CollectionPolicyUpdaterHint, fact, fs, memo)
	if err != nil {
		return CollectionPolicyUpdater{}, err
	}
	return CollectionPolicyUpdater{BaseOperation: bo}, nil
}
