package collection

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"strconv"
)

var MaxApproveItems = 10

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
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
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

type Approve struct {
	currency.BaseOperation
}

func NewApprove(fact ApproveFact) (Approve, error) {
	return Approve{BaseOperation: currency.NewBaseOperation(ApproveHint, fact)}, nil
}

func (op *Approve) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
