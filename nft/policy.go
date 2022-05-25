package nft

import (
	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type PaymentParameter uint

func (pp PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(pp))
}

func (pp PaymentParameter) Uint() uint {
	return uint(pp)
}

func (pp PaymentParameter) IsValid([]byte) error {
	if uint(pp) >= 100 {
		return isvalid.InvalidError.Errorf(
			"invalid range of symbol; %d <= %d < %d", 0, pp, 100)
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
	symbol  extension.ContractID
	policy  BasePolicy
}

func NewDesign(parent base.Address, creator base.Address, symbol extension.ContractID, policy BasePolicy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		parent:     parent,
		creator:    creator,
		symbol:     symbol,
		policy:     policy,
	}
}

func MustNewDesign(parent base.Address, creator base.Address, symbol extension.ContractID, policy BasePolicy) Design {
	d := NewDesign(parent, creator, symbol, policy)
	if err := d.IsValid(nil); err != nil {
		panic(err)
	}
	return d
}

func (d Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.parent.Bytes(),
		d.creator.Bytes(),
		d.symbol.Bytes(),
		d.policy.Bytes(),
	)
}

func (d Design) Parent() base.Address {
	return d.parent
}

func (d Design) Creator() base.Address {
	return d.creator
}

func (d Design) Symbol() extension.ContractID {
	return d.symbol
}

func (d Design) Policy() BasePolicy {
	return d.policy
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
	return nil
}

func (d Design) Rebuild() Design {
	d.policy = d.policy.Rebuild()
	return d
}

type BasePolicy interface {
	isvalid.IsValider
	Bytes() []byte
	Rebuild() BasePolicy
}
