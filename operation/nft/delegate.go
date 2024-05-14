package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	DelegateFactHint = hint.MustNewHint("mitum-nft-delegate-operation-fact-v0.0.1")
	DelegateHint     = hint.MustNewHint("mitum-nft-delegate-operation-v0.0.1")
)

var MaxDelegateItems = 100

type DelegateFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []DelegateItem
}

func NewDelegateFact(token []byte, sender mitumbase.Address, items []DelegateItem) DelegateFact {
	bf := mitumbase.NewBaseFact(DelegateFactHint, token)
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
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(errors.Errorf("empty items for DelegateFact"))
	} else if l > int(MaxDelegateItems) {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxDelegateItems)))
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(errors.Errorf("sender is same with contract"))
		}

		delegatee := item.Delegatee()

		if addressMap, collectionFound := founds[item.contract.String()]; !collectionFound {
			founds[item.contract.String()] = make(map[string]struct{})
		} else if _, addressFound := addressMap[delegatee.String()]; addressFound {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("collection-operator, %v, %v", item.contract, delegatee)))
		}

		founds[item.contract.String()][delegatee.String()] = struct{}{}
	}
	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
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

func (fact DelegateFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact DelegateFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact DelegateFact) Addresses() ([]mitumbase.Address, error) {
	l := len(fact.items)

	as := make([]mitumbase.Address, l+1)

	for i, item := range fact.items {
		as[i] = item.Delegatee()
	}

	as[l] = fact.sender

	return as, nil
}

func (fact DelegateFact) Items() []DelegateItem {
	return fact.items
}

type Delegate struct {
	common.BaseOperation
}

func NewDelegate(fact DelegateFact) (Delegate, error) {
	return Delegate{BaseOperation: common.NewBaseOperation(DelegateHint, fact)}, nil
}
