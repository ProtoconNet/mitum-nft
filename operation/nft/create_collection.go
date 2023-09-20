package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateCollectionFactHint = hint.MustNewHint("mitum-nft-create-collection-operation-fact-v0.0.1")
	CreateCollectionHint     = hint.MustNewHint("mitum-nft-create-collection-operation-v0.0.1")
)

type CreateCollectionFact struct {
	mitumbase.BaseFact
	sender    mitumbase.Address
	contract  mitumbase.Address
	name      types.CollectionName
	royalty   types.PaymentParameter
	uri       types.URI
	whitelist []mitumbase.Address
	currency  currencytypes.CurrencyID
}

func NewCreateCollectionFact(
	token []byte,
	sender mitumbase.Address,
	contract mitumbase.Address,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []mitumbase.Address,
	currency currencytypes.CurrencyID,
) CreateCollectionFact {
	bf := mitumbase.NewBaseFact(CreateCollectionFactHint, token)
	fact := CreateCollectionFact{
		BaseFact:  bf,
		sender:    sender,
		contract:  contract,
		name:      name,
		royalty:   royalty,
		uri:       uri,
		whitelist: whitelist,
		currency:  currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateCollectionFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.name,
		fact.royalty,
		fact.uri,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("sender and contract are the same, %q == %q", fact.sender, fact.contract)
	}

	return nil
}

func (fact CreateCollectionFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateCollectionFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateCollectionFact) Bytes() []byte {
	as := make([][]byte, len(fact.whitelist))
	for i, white := range fact.whitelist {
		as[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.name.Bytes(),
		fact.royalty.Bytes(),
		fact.uri.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (fact CreateCollectionFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CreateCollectionFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CreateCollectionFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact CreateCollectionFact) Name() types.CollectionName {
	return fact.name
}

func (fact CreateCollectionFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact CreateCollectionFact) URI() types.URI {
	return fact.uri
}

func (fact CreateCollectionFact) WhiteList() []mitumbase.Address {
	return fact.whitelist
}

func (fact CreateCollectionFact) Addresses() ([]mitumbase.Address, error) {
	l := 2 + len(fact.whitelist)

	as := make([]mitumbase.Address, l)
	copy(as, fact.whitelist)

	as[l-2] = fact.sender
	as[l-1] = fact.contract

	return as, nil
}

func (fact CreateCollectionFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

type CreateCollection struct {
	common.BaseOperation
}

func NewCreateCollection(fact CreateCollectionFact) (CreateCollection, error) {
	return CreateCollection{BaseOperation: common.NewBaseOperation(CreateCollectionHint, fact)}, nil
}

func (op *CreateCollection) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
