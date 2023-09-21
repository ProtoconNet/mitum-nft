package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxApproveItems = 10

var (
	ApproveFactHint = hint.MustNewHint("mitum-nft-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-nft-approve-operation-v0.0.1")
)

type ApproveFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []ApproveItem
}

func NewApproveFact(token []byte, sender mitumbase.Address, items []ApproveItem) ApproveFact {
	bf := mitumbase.NewBaseFact(ApproveFactHint, token)
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
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for ApproveFact")
	} else if l > int(MaxApproveItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxApproveItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return err
		}

		n := strconv.FormatUint(item.NFT(), 10)

		if _, found := founds[n]; found {
			return util.ErrInvalid.Errorf("duplicate nft found, %q", n)
		}

		founds[n] = struct{}{}
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

func (fact ApproveFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact ApproveFact) Items() []ApproveItem {
	return fact.items
}

func (fact ApproveFact) Addresses() ([]mitumbase.Address, error) {
	as := make([]mitumbase.Address, len(fact.items)+1)

	for i := range fact.items {
		as[i] = fact.items[i].Approved()
	}
	as[len(fact.items)] = fact.sender

	return as, nil
}

type Approve struct {
	common.BaseOperation
}

func NewApprove(fact ApproveFact) (Approve, error) {
	return Approve{BaseOperation: common.NewBaseOperation(ApproveHint, fact)}, nil
}
