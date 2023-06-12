package types

import (
	"bytes"
	"regexp"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MaxWhites = 10

var (
	MinLengthCollectionName = 3
	MaxLengthCollectionName = 30
	ReValidCollectionName   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\s]+$`)
)

type CollectionName string

func (cn CollectionName) IsValid([]byte) error {
	l := len(cn)

	if l < MinLengthCollectionName {
		return util.ErrInvalid.Errorf(
			"collection name length under min, %d < %d", l, MinLengthCollectionName)
	}

	if l > MaxLengthCollectionName {
		return util.ErrInvalid.Errorf(
			"collection name length over max, %d > %d", l, MaxLengthCollectionName)
	}

	if !ReValidCollectionName.Match([]byte(cn)) {
		return util.ErrInvalid.Errorf("wrong collection name, %q", cn)
	}

	return nil
}

func (cn CollectionName) Bytes() []byte {
	return []byte(cn)
}

func (cn CollectionName) String() string {
	return string(cn)
}

var CollectionPolicyHint = hint.MustNewHint("mitum-nft-collection-policy-v0.0.1")

type CollectionPolicy struct {
	hint.BaseHinter
	name    CollectionName
	royalty PaymentParameter
	uri     URI
	whites  []mitumbase.Address
}

func NewCollectionPolicy(name CollectionName, royalty PaymentParameter, uri URI, whites []mitumbase.Address) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whites:     whites,
	}
}

func (policy CollectionPolicy) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		policy.name,
		policy.royalty,
		policy.uri,
	); err != nil {
		return err
	}

	if l := len(policy.whites); l > MaxWhites {
		return util.ErrInvalid.Errorf("whites over allowed, %d > %d", l, MaxWhites)
	}

	founds := map[string]struct{}{}
	for _, white := range policy.whites {
		if err := white.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[white.String()]; found {
			return util.ErrInvalid.Errorf("duplicate white found, %q", white)
		}
		founds[white.String()] = struct{}{}
	}

	return nil
}

func (policy CollectionPolicy) Bytes() []byte {
	as := make([][]byte, len(policy.whites))
	for i, white := range policy.whites {
		as[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		policy.name.Bytes(),
		policy.royalty.Bytes(),
		policy.uri.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (policy CollectionPolicy) Name() CollectionName {
	return policy.name
}

func (policy CollectionPolicy) Royalty() PaymentParameter {
	return policy.royalty
}

func (policy CollectionPolicy) URI() URI {
	return policy.uri
}

func (policy CollectionPolicy) Whites() []mitumbase.Address {
	return policy.whites
}

func (policy CollectionPolicy) Addresses() ([]mitumbase.Address, error) {
	return policy.whites, nil
}

func (policy CollectionPolicy) Equal(c BasePolicy) bool {
	cpolicy, ok := c.(CollectionPolicy)
	if !ok {
		return false
	}

	if policy.name != cpolicy.name {
		return false
	}

	if policy.royalty != cpolicy.royalty {
		return false
	}

	if policy.uri != cpolicy.uri {
		return false
	}

	if len(policy.whites) != len(cpolicy.whites) {
		return false
	}

	whites := policy.Whites()
	cwhites := cpolicy.Whites()
	sort.Slice(whites, func(i, j int) bool {
		return bytes.Compare(whites[j].Bytes(), whites[i].Bytes()) < 0
	})
	sort.Slice(cwhites, func(i, j int) bool {
		return bytes.Compare(cwhites[j].Bytes(), cwhites[i].Bytes()) < 0
	})

	for i := range whites {
		if !whites[i].Equal(cwhites[i]) {
			return false
		}
	}

	return true
}

var CollectionDesignHint = hint.MustNewHint("mitum-nft-collection-design-v0.0.1")

type CollectionDesign struct {
	Design
}

func NewCollectionDesign(parent mitumbase.Address, creator mitumbase.Address, collection types.ContractID, active bool, policy CollectionPolicy) CollectionDesign {
	design := NewDesign(parent, creator, collection, active, policy)

	return CollectionDesign{
		Design: design,
	}
}
