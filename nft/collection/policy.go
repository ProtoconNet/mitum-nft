package collection

import (
	"bytes"
	"regexp"
	"sort"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxWhiteAddress = 10

var (
	MinLengthCollectionName = 3
	MaxLengthCollectionName = 30
	ReValidCollectionName   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\s]+$`)
)

type CollectionName string

func (cn CollectionName) Bytes() []byte {
	return []byte(cn)
}

func (cn CollectionName) String() string {
	return string(cn)
}

func (cn CollectionName) IsValid([]byte) error {
	if l := len(cn); l < MinLengthCollectionName || l > MaxLengthCollectionName {
		return isvalid.InvalidError.Errorf(
			"invalid length of collection name; %d <= %d <= %d", MinLengthCollectionName, l, MaxLengthCollectionName)
	} else if !ReValidCollectionName.Match([]byte(cn)) {
		return isvalid.InvalidError.Errorf("wrong collection name; %q", cn)
	}

	return nil
}

var (
	CollectionPolicyType   = hint.Type("mitum-nft-collection-policy")
	CollectionPolicyHint   = hint.NewHint(CollectionPolicyType, "v0.0.1")
	CollectionPolicyHinter = CollectionPolicy{BaseHinter: hint.NewBaseHinter(CollectionPolicyHint)}
)

type CollectionPolicy struct {
	hint.BaseHinter
	name    CollectionName
	royalty nft.PaymentParameter
	uri     nft.URI
	whites  []base.Address
}

func NewCollectionPolicy(name CollectionName, royalty nft.PaymentParameter, uri nft.URI, whites []base.Address) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whites:     whites,
	}
}

func MustNewCollectionPolicy(name CollectionName, royalty nft.PaymentParameter, uri nft.URI, whites []base.Address) CollectionPolicy {
	policy := NewCollectionPolicy(name, royalty, uri, whites)

	if err := policy.IsValid(nil); err != nil {
		panic(err)
	}

	return policy
}

func (policy CollectionPolicy) Bytes() []byte {
	as := make([][]byte, len(policy.whites))
	for i := range policy.whites {
		as[i] = policy.whites[i].Bytes()
	}

	return util.ConcatBytesSlice(
		policy.name.Bytes(),
		policy.royalty.Bytes(),
		policy.uri.Bytes(),
		util.ConcatBytesSlice(as...),
	)
}

func (policy CollectionPolicy) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		policy.name,
		policy.royalty,
		policy.uri); err != nil {
		return err
	}

	if l := len(policy.whites); l > MaxWhiteAddress {
		return isvalid.InvalidError.Errorf("address in white list over allowed; %d > %d", l, MaxWhiteAddress)
	}

	founds := map[base.Address]struct{}{}
	for i := range policy.whites {
		acc := policy.whites[i]
		if err := acc.IsValid(nil); err != nil {
			return err
		}
		if _, found := founds[acc]; found {
			return isvalid.InvalidError.Errorf("duplicate white found; %q", acc)
		}
		founds[acc] = struct{}{}
	}

	return nil
}

func (policy CollectionPolicy) Name() CollectionName {
	return policy.name
}

func (policy CollectionPolicy) Royalty() nft.PaymentParameter {
	return policy.royalty
}

func (policy CollectionPolicy) Uri() nft.URI {
	return policy.uri
}

func (policy CollectionPolicy) Whites() []base.Address {
	return policy.whites
}

func (policy CollectionPolicy) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(policy.whites))
	for i := range policy.whites {
		as[i] = policy.whites[i]
	}
	return as, nil
}

func (policy CollectionPolicy) Equal(c nft.BasePolicy) bool {
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

func (policy CollectionPolicy) Rebuild() nft.BasePolicy {
	return policy
}
