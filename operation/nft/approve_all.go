package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	ApproveAllFactHint = hint.MustNewHint("mitum-nft-approve-all-operation-fact-v0.0.1")
	ApproveAllHint     = hint.MustNewHint("mitum-nft-approve-all-operation-v0.0.1")
)

var MaxApproveAllItems = 100

type ApproveAllFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []ApproveAllItem
}

func NewApproveAllFact(token []byte, sender mitumbase.Address, items []ApproveAllItem) ApproveAllFact {
	bf := mitumbase.NewBaseFact(ApproveAllFactHint, token)
	fact := ApproveAllFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ApproveAllFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(errors.Errorf("empty items for DelegateFact"))
	} else if l > int(MaxApproveAllItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxApproveAllItems)))
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
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(
				errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		if fact.sender.Equal(item.approved) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(
				errors.Errorf("sender %v delegates to itself", fact.sender)))
		}

		if addressMap, collectionFound := founds[item.contract.String()]; !collectionFound {
			founds[item.contract.String()] = make(map[string]struct{})
		} else if _, addressFound := addressMap[item.Approved().String()]; addressFound {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf("approved %v in contract account %v", item.Approved(), item.contract)))
		}

		founds[item.contract.String()][item.Approved().String()] = struct{}{}
	}
	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact ApproveAllFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ApproveAllFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveAllFact) Bytes() []byte {
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

func (fact ApproveAllFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveAllFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact ApproveAllFact) Addresses() ([]mitumbase.Address, error) {
	l := len(fact.items)

	as := make([]mitumbase.Address, l+1)

	for i, item := range fact.items {
		as[i] = item.Approved()
	}

	as[l] = fact.sender

	return as, nil
}

func (fact ApproveAllFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)

	for i := range fact.items {
		zeroBig := common.ZeroBig
		cid := fact.items[i].Currency()
		var amsTemp []common.Big
		if ams, found := required[cid]; found {
			ams = append(ams, zeroBig)
			required[cid] = ams
		} else {
			amsTemp = append(amsTemp, zeroBig)
			required[cid] = amsTemp
		}
	}

	return required
}

func (fact ApproveAllFact) FeePayer() mitumbase.Address {
	return fact.sender
}

func (fact ApproveAllFact) FactUser() mitumbase.Address {
	return fact.sender
}

func (fact ApproveAllFact) ActiveContract() []mitumbase.Address {
	var arr []mitumbase.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact ApproveAllFact) Items() []ApproveAllItem {
	return fact.items
}

type ApproveAll struct {
	extras.ExtendedOperation
}

func NewDelegate(fact ApproveAllFact) (ApproveAll, error) {
	return ApproveAll{
		ExtendedOperation: extras.NewExtendedOperation(ApproveAllHint, fact),
	}, nil
}
