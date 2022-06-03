package collection

import (
	"regexp"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MinLengthCollectionName = 3
	MaxLengthCollectionName = 30
	ReValidCollectionName   = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
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
}

func NewCollectionPolicy(name CollectionName, royalty nft.PaymentParameter, uri nft.URI) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		name:       name,
		royalty:    royalty,
		uri:        uri,
	}
}

func MustNewCollectionPolicy(name CollectionName, royalty nft.PaymentParameter, uri nft.URI) CollectionPolicy {
	policy := NewCollectionPolicy(name, royalty, uri)

	if err := policy.IsValid(nil); err != nil {
		panic(err)
	}

	return policy
}

func (policy CollectionPolicy) Bytes() []byte {
	return util.ConcatBytesSlice(
		policy.name.Bytes(),
		policy.royalty.Bytes(),
		policy.uri.Bytes(),
	)
}

func (policy CollectionPolicy) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		policy.name,
		policy.royalty); err != nil {
		return err
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

func (policy CollectionPolicy) Addresses() ([]base.Address, error) {
	as := []base.Address{}
	return as, nil
}

func (policy CollectionPolicy) Rebuild() nft.BasePolicy {
	return policy
}
