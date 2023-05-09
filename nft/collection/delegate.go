package collection

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	DelegateFactHint = hint.MustNewHint("mitum-nft-delegate-operation-fact-v0.0.1")
	DelegateHint     = hint.MustNewHint("mitum-nft-delegate-operation-v0.0.1")
)

var (
	MaxOperators     = 10
	MaxDelegateItems = 10
)

type DelegateFact struct {
	base.BaseFact
	sender base.Address
	items  []DelegateItem
}

func NewDelegateFact(token []byte, sender base.Address, items []DelegateItem) DelegateFact {
	bf := base.NewBaseFact(DelegateFactHint, token)
	fact := DelegateFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact DelegateFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for DelegateFact")
	} else if l > int(MaxDelegateItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxDelegateItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return err
		}

		operator := item.Operator()
		collection := item.Collection()

		if addressMap, collectionFound := founds[collection.String()]; !collectionFound {
			founds[collection.String()] = make(map[string]struct{})
		} else if _, addressFound := addressMap[operator.String()]; addressFound {
			return util.ErrInvalid.Errorf("duplicate collection-operator found, %q-%q", collection, operator)
		}

		founds[collection.String()][operator.String()] = struct{}{}
	}

	return nil
}

func (fact DelegateFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact DelegateFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DelegateFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i, item := range fact.items {
		is[i] = item.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact DelegateFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact DelegateFact) Sender() base.Address {
	return fact.sender
}

func (fact DelegateFact) Addresses() ([]base.Address, error) {
	l := len(fact.items)

	as := make([]base.Address, l+1)

	for i, item := range fact.items {
		as[i] = item.Operator()
	}

	as[l] = fact.sender

	return as, nil
}

func (fact DelegateFact) Items() []DelegateItem {
	return fact.items
}

type Delegate struct {
	currency.BaseOperation
}

func NewDelegate(fact DelegateFact) (Delegate, error) {
	return Delegate{BaseOperation: currency.NewBaseOperation(DelegateHint, fact)}, nil
}

func (op *Delegate) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
