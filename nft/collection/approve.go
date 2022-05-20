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
	h        valuehash.Hash
	token    []byte
	sender   base.Address
	approved base.Address
	nfts     []nft.NFTID
	cid      currency.CurrencyID
}

func NewApproveFact(token []byte, sender base.Address, approved base.Address, nfts []nft.NFTID, cid currency.CurrencyID) ApproveFact {
	fact := ApproveFact{
		BaseHinter: hint.NewBaseHinter(ApproveFactHint),
		token:      token,
		sender:     sender,
		approved:   approved,
		nfts:       nfts,
		cid:        cid,
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
	ns := make([][]byte, len(fact.nfts))
	for i := range fact.nfts {
		ns[i] = fact.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.approved.Bytes(),
		fact.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for ApproveFact")
	} else if len(fact.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for ApproveFact")
	}

	if err := isvalid.Check(
		nil, false, fact.h,
		fact.sender, fact.approved, fact.cid); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.nfts {
		if err := fact.nfts[i].IsValid(nil); err != nil {
			return err
		}
		nft := fact.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found; %s", nft)
		}
		foundNFT[nft] = true
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

func (fact ApproveFact) Approved() base.Address {
	return fact.approved
}

func (fact ApproveFact) NFTs() []nft.NFTID {
	return fact.nfts
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.approved

	return as, nil
}

func (fact ApproveFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact ApproveFact) Rebuild() ApproveFact {
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
