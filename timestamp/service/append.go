package service

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxAppendItems = 10

var (
	AppendFactHint = hint.MustNewHint("mitum-timestamp-append-operation-fact-v0.0.1")
	AppendHint     = hint.MustNewHint("mitum-timestamp-append-operation-v0.0.1")
)

type AppendFact struct {
	base.BaseFact
	sender           base.Address
	target           base.Address
	service          extensioncurrency.ContractID
	projectID        string
	requestTimeStamp uint64
	data             string
	currency         currency.CurrencyID
}

func NewAppendFact(token []byte, sender, target base.Address, service extensioncurrency.ContractID, projectID string, requestTimeStamp uint64, data string, currency currency.CurrencyID) AppendFact {
	bf := base.NewBaseFact(AppendFactHint, token)
	fact := AppendFact{
		BaseFact:         bf,
		sender:           sender,
		target:           target,
		service:          service,
		projectID:        projectID,
		requestTimeStamp: requestTimeStamp,
		data:             data,
		currency:         currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact AppendFact) IsValid(b []byte) error {
	if len(fact.projectID) < 1 || len(fact.projectID) > timestamp.MaxProjectIDLen {
		return errors.Errorf("invalid projectID length %v < 1 or > %v", len(fact.projectID), timestamp.MaxProjectIDLen)
	}

	if len(fact.data) < 1 || len(fact.data) > timestamp.MaxDataLen {
		return errors.Errorf("invalid data length %v < 1 or > %v", len(fact.data), timestamp.MaxDataLen)
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.target,
		fact.service,
		fact.currency,
	); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact AppendFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AppendFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AppendFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.service.Bytes(),
		[]byte(fact.projectID),
		util.Uint64ToBytes(fact.requestTimeStamp),
		[]byte(fact.data),
		fact.currency.Bytes(),
	)
}

func (fact AppendFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AppendFact) Sender() base.Address {
	return fact.sender
}

func (fact AppendFact) Target() base.Address {
	return fact.target
}

type Append struct {
	currency.BaseOperation
}

func NewAppend(fact AppendFact) (Append, error) {
	return Append{BaseOperation: currency.NewBaseOperation(AppendHint, fact)}, nil
}

func (op *Append) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
