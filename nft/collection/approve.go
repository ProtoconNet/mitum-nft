package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxApproveItems = 10

var (
	ApproveFactType   = hint.Type("mitum-nft-approve-operation-fact")
	ApproveFactHint   = hint.NewHint(ApproveFactType, "v0.0.1")
	ApproveFactHinter = ApproveFact{BaseHinter: hint.NewBaseHinter(ApproveFactHint)}
	ApproveType       = hint.Type("mitum-nft-approve-operation")
	ApproveHint       = hint.NewHint(ApproveType, "v0.0.1")
	ApproveHinter     = Approve{BaseOperation: operationHinter(ApproveHint)}
)

type ApproveFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []ApproveItem
}

func NewApproveFact(token []byte, sender base.Address, items []ApproveItem) ApproveFact {
	fact := ApproveFact{
		BaseHinter: hint.NewBaseHinter(ApproveFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact ApproveFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact ApproveFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return isvalid.InvalidError.Errorf("empty items for ApproveFact")
	} else if l > int(MaxApproveItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", l, MaxApproveItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[nft.NFTID]struct{}{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		n := fact.items[i].NFT()
		if err := n.IsValid(nil); err != nil {
			return err
		}

		if _, found := founds[n]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found; %q", n)
		}

		founds[n] = struct{}{}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact ApproveFact) Token() []byte {
	return fact.token
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) Items() []ApproveItem {
	return fact.items
}

func (fact ApproveFact) NFTs() []nft.NFTID {
	ns := make([]nft.NFTID, len(fact.items))

	for i := range fact.items {
		ns[i] = fact.items[i].NFT()
	}

	return ns
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	for i := range fact.items {
		as[i] = fact.items[i].Approved()
	}
	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact ApproveFact) Currencies() []currency.CurrencyID {
	cs := make([]currency.CurrencyID, len(fact.items))

	for i := range fact.items {
		cs[i] = fact.items[i].Currency()
	}

	return cs
}

func (fact ApproveFact) Rebuild() ApproveFact {
	items := make([]ApproveItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Approve struct {
	currency.BaseOperation
}

func NewApprove(fact ApproveFact, fs []base.FactSign, memo string) (Approve, error) {
	bo, err := currency.NewBaseOperationFromFact(ApproveHint, fact, fs, memo)
	if err != nil {
		return Approve{}, err
	}
	return Approve{BaseOperation: bo}, nil
}
