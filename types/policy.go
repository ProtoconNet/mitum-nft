package types

import (
	"bytes"
	"regexp"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/common"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var MaxWhitelist = 20

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
		return util.ErrInvalid.Errorf("wrong collection name, %v", cn)
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
	name      CollectionName
	royalty   PaymentParameter
	uri       URI
	whitelist []mitumbase.Address
}

func NewCollectionPolicy(name CollectionName, royalty PaymentParameter, uri URI, whitelist []mitumbase.Address) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whitelist:  whitelist,
	}
}

func (policy CollectionPolicy) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		policy.BaseHinter,
		policy.name,
		policy.royalty,
		policy.uri,
	); err != nil {
		return err
	}

	if l := len(policy.whitelist); l > MaxWhitelist {
		return common.ErrArrayLen.Wrap(errors.Errorf("whitelist over allowed, %d > %d", l, MaxWhitelist))
	}

	founds := map[string]struct{}{}
	for _, white := range policy.whitelist {
		if err := white.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[white.String()]; found {
			return common.ErrDupVal.Wrap(errors.Errorf("duplicate white found, %v", white))
		}
		founds[white.String()] = struct{}{}
	}

	return nil
}

func (policy CollectionPolicy) Bytes() []byte {
	as := make([][]byte, len(policy.whitelist))
	for i, white := range policy.whitelist {
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

func (policy CollectionPolicy) Whitelist() []mitumbase.Address {
	return policy.whitelist
}

func (policy CollectionPolicy) Addresses() ([]mitumbase.Address, error) {
	return policy.whitelist, nil
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

	if len(policy.whitelist) != len(cpolicy.whitelist) {
		return false
	}

	whitelist := policy.Whitelist()
	cwhitelist := cpolicy.Whitelist()
	sort.Slice(whitelist, func(i, j int) bool {
		return bytes.Compare(whitelist[j].Bytes(), whitelist[i].Bytes()) < 0
	})
	sort.Slice(cwhitelist, func(i, j int) bool {
		return bytes.Compare(cwhitelist[j].Bytes(), cwhitelist[i].Bytes()) < 0
	})

	for i := range whitelist {
		if !whitelist[i].Equal(cwhitelist[i]) {
			return false
		}
	}

	return true
}

var CollectionDesignHint = hint.MustNewHint("mitum-nft-collection-design-v0.0.1")

type CollectionDesign struct {
	Design
}

func NewCollectionDesign(contract mitumbase.Address, creator mitumbase.Address, active bool, policy CollectionPolicy) CollectionDesign {
	design := NewDesign(contract, creator, active, policy)

	return CollectionDesign{
		Design: design,
	}
}
