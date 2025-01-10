package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxApproveItems = 100

var (
	ApproveFactHint = hint.MustNewHint("mitum-nft-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-nft-approve-operation-v0.0.1")
)

type ApproveFact struct {
	base.BaseFact
	sender base.Address
	items  []ApproveItem
}

func NewApproveFact(token []byte, sender base.Address, items []ApproveItem) ApproveFact {
	bf := base.NewBaseFact(ApproveFactHint, token)
	fact := ApproveFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}

	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if n := len(fact.items); n < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items")))
	} else if n > int(MaxApproveItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items, %d over max, %d", n, MaxApproveItems)))
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		n := strconv.FormatUint(item.NFTIdx(), 10)

		if _, found := founds[n]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("nft idx %v", n)))
		}

		founds[n] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact ApproveFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
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

func (fact ApproveFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) Items() []ApproveItem {
	return fact.items
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	for i := range fact.items {
		as[i] = fact.items[i].Approved()
	}
	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact ApproveFact) FeeBase() map[types.CurrencyID][]common.Big {
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

func (fact ApproveFact) FeePayer() base.Address {
	return fact.sender
}

func (fact ApproveFact) FactUser() base.Address {
	return fact.sender
}

func (fact ApproveFact) Signer() base.Address {
	return fact.sender
}

func (fact ApproveFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

type Approve struct {
	extras.ExtendedOperation
}

func NewApprove(fact ApproveFact) (Approve, error) {
	return Approve{
		ExtendedOperation: extras.NewExtendedOperation(ApproveHint, fact),
	}, nil
}
