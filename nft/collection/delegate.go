package collection

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	DelegateFactType   = hint.Type("mitum-nft-delegate-operation-fact")
	DelegateFactHint   = hint.NewHint(DelegateFactType, "v0.0.1")
	DelegateFactHinter = DelegateFact{BaseHinter: hint.NewBaseHinter(DelegateFactHint)}
	DelegateType       = hint.Type("mitum-nft-delegate-operation")
	DelegateHint       = hint.NewHint(DelegateType, "v0.0.1")
	DelegateHinter     = Delegate{BaseOperation: operationHinter(DelegateHint)}
)

var MaxAgents = 10

var (
	DelegateAllow  = DelegateMode("allow")
	DelegateCancel = DelegateMode("cancel")
)

type DelegateMode string

func (mode DelegateMode) Bytes() []byte {
	return []byte(mode)
}

func (mode DelegateMode) String() string {
	return string(mode)
}

func (mode DelegateMode) IsValid([]byte) error {
	if !(mode == DelegateAllow || mode == DelegateCancel) {
		return isvalid.InvalidError.Errorf("wrong delegate mode; %s", mode)
	}

	return nil
}

type DelegateFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	agents []base.Address
	mode   DelegateMode
	cid    currency.CurrencyID
}

func NewDelegateFact(token []byte, sender base.Address, agents []base.Address, mode DelegateMode, cid currency.CurrencyID) DelegateFact {
	fact := DelegateFact{
		BaseHinter: hint.NewBaseHinter(DelegateFactHint),
		token:      token,
		sender:     sender,
		agents:     agents,
		mode:       mode,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact DelegateFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact DelegateFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DelegateFact) Bytes() []byte {
	ags := make([][]byte, len(fact.agents))

	for i := range fact.agents {
		ags[i] = fact.agents[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.mode.Bytes(),
		fact.cid.Bytes(),
		util.ConcatBytesSlice(ags...),
	)
}

func (fact DelegateFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for DelegateFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.mode,
		fact.cid); err != nil {
		return err
	}

	if n := len(fact.agents); n < 1 {
		return isvalid.InvalidError.Errorf("empty agents for DelegateFact")
	} else if n > MaxAgents {
		return isvalid.InvalidError.Errorf("agents over allowed; %d > %d", n, MaxAgents)
	}

	foundAgent := map[string]bool{}
	for i := range fact.agents {
		if err := fact.agents[i].IsValid(nil); err != nil {
			return err
		}

		agent := fact.agents[i].String()
		if _, found := foundAgent[agent]; found {
			return isvalid.InvalidError.Errorf("duplicate agent found; %s", agent)
		}

		foundAgent[agent] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact DelegateFact) Token() []byte {
	return fact.token
}

func (fact DelegateFact) Sender() base.Address {
	return fact.sender
}

func (fact DelegateFact) Agents() []base.Address {
	return fact.agents
}

func (fact DelegateFact) Mode() DelegateMode {
	return fact.mode
}

func (fact DelegateFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.agents)+1)

	for i := range fact.agents {
		as[i] = fact.agents[i]
	}

	as[len(fact.agents)] = fact.sender

	return as, nil
}

func (fact DelegateFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact DelegateFact) Rebuild() DelegateFact {
	fact.h = fact.GenerateHash()

	return fact
}

type Delegate struct {
	currency.BaseOperation
}

func NewDelegate(fact DelegateFact, fs []base.FactSign, memo string) (Delegate, error) {
	bo, err := currency.NewBaseOperationFromFact(DelegateHint, fact, fs, memo)
	if err != nil {
		return Delegate{}, err
	}
	return Delegate{BaseOperation: bo}, nil
}
