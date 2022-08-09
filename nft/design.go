package nft

import (
	"net/url"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxPaymentParameter uint = 99

type PaymentParameter uint

func (pp PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(pp))
}

func (pp PaymentParameter) Uint() uint {
	return uint(pp)
}

func (pp PaymentParameter) IsValid([]byte) error {
	if uint(pp) > MaxPaymentParameter {
		return isvalid.InvalidError.Errorf(
			"invalid range of paymentparameter; %d <= %d <= %d", 0, pp, MaxPaymentParameter)
	}

	return nil
}

var MaxURILength = 1000

type URI string

func (uri URI) Bytes() []byte {
	return []byte(uri)
}

func (uri URI) String() string {
	return string(uri)
}

func (uri URI) IsValid([]byte) error {
	if _, err := url.Parse(string(uri)); err != nil {
		return err
	}

	if l := len(uri); l > 1000 {
		return isvalid.InvalidError.Errorf("invalid length of uri; %d > %d", l, MaxURILength)
	}

	if uri != "" && strings.TrimSpace(string(uri)) == "" {
		return isvalid.InvalidError.Errorf("uri with only spaces")
	}

	return nil
}

var (
	DesignType   = hint.Type("mitum-nft-design")
	DesignHint   = hint.NewHint(DesignType, "v0.0.1")
	DesignHinter = Design{BaseHinter: hint.NewBaseHinter(DesignHint)}
)

type Design struct {
	hint.BaseHinter
	parent  base.Address
	creator base.Address
	symbol  extensioncurrency.ContractID
	active  bool
	policy  BasePolicy
}

func NewDesign(parent base.Address, creator base.Address, symbol extensioncurrency.ContractID, active bool, policy BasePolicy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		parent:     parent,
		creator:    creator,
		symbol:     symbol,
		active:     active,
		policy:     policy,
	}
}

func MustNewDesign(parent base.Address, creator base.Address, symbol extensioncurrency.ContractID, active bool, policy BasePolicy) Design {
	d := NewDesign(parent, creator, symbol, active, policy)
	if err := d.IsValid(nil); err != nil {
		panic(err)
	}
	return d
}

func (d Design) Bytes() []byte {
	ab := make([]byte, 1)
	if d.active {
		ab[0] = 1
	} else {
		ab[0] = 0
	}

	return util.ConcatBytesSlice(
		d.parent.Bytes(),
		d.creator.Bytes(),
		d.symbol.Bytes(),
		ab,
		d.policy.Bytes(),
	)
}

func (d Design) Hint() hint.Hint {
	return DesignHint
}

func (d Design) Hash() valuehash.Hash {
	return d.GenerateHash()
}

func (d Design) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(d.Bytes())
}

func (d Design) Parent() base.Address {
	return d.parent
}

func (d Design) Creator() base.Address {
	return d.creator
}

func (d Design) Symbol() extensioncurrency.ContractID {
	return d.symbol
}

func (d Design) Active() bool {
	return d.active
}

func (d Design) Policy() BasePolicy {
	return d.policy
}

func (d Design) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = d.parent
	as[1] = d.creator

	if ads, err := d.Policy().Addresses(); err != nil {
		return as, err
	} else {
		as = append(as, ads...)
	}

	return as, nil
}

func (d Design) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		d.BaseHinter,
		d.parent,
		d.creator,
		d.symbol,
		d.policy); err != nil {
		return err
	}

	if d.parent.Equal(d.creator) {
		return isvalid.InvalidError.Errorf("parent and creator are the same; %q == %q", d.parent, d.creator)
	}

	return nil
}

func (d Design) Equal(cd Design) bool {
	if !d.parent.Equal(cd.parent) {
		return false
	}

	if !d.creator.Equal(cd.creator) {
		return false
	}

	if d.symbol != cd.symbol {
		return false
	}

	if d.active != cd.active {
		return false
	}

	if !d.policy.Equal(cd.policy) {
		return false
	}

	if d.Hash() != cd.Hash() {
		return false
	}

	return true
}

func (d Design) Rebuild() Design {
	d.policy = d.policy.Rebuild()
	return d
}

type BasePolicy interface {
	isvalid.IsValider
	Bytes() []byte
	Addresses() ([]base.Address, error)
	Equal(c BasePolicy) bool
	Rebuild() BasePolicy
}
