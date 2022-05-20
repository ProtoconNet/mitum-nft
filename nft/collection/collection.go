package collection

import (
	"bytes"
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

type CollectionUri string

func (cu CollectionUri) Bytes() []byte {
	return []byte(cu)
}

func (cu CollectionUri) String() string {
	return string(cu)
}

func (cu CollectionUri) IsValid([]byte) error {
	return nil
}

var (
	CollectionPolicyType   = hint.Type("mitum-nft-collection-policy")
	CollectionPolicyHint   = hint.NewHint(CollectionPolicyType, "v0.0.1")
	CollectionPolicyHinter = CollectionPolicy{BaseHinter: hint.NewBaseHinter(CollectionPolicyHint)}
)

type CollectionPolicy struct {
	hint.BaseHinter
	symbol  nft.Symbol
	name    CollectionName
	creator base.Address
	royalty nft.PaymentParameter
	uri     CollectionUri
}

func NewCollectionPolicy(symbol nft.Symbol, name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri CollectionUri) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		symbol:     symbol,
		name:       name,
		creator:    creator,
		royalty:    royalty,
		uri:        uri,
	}
}

func MustNewCollectionPolicy(symbol nft.Symbol, name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri CollectionUri) CollectionPolicy {
	policy := NewCollectionPolicy(symbol, name, creator, royalty, uri)

	if err := policy.IsValid(nil); err != nil {
		panic(err)
	}

	return policy
}

func (policy CollectionPolicy) Bytes() []byte {
	return util.ConcatBytesSlice(
		policy.symbol.Bytes(),
		policy.name.Bytes(),
		policy.creator.Bytes(),
		policy.royalty.Bytes(),
		policy.uri.Bytes(),
	)
}

func (policy CollectionPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		policy.symbol,
		policy.name,
		policy.creator,
		policy.royalty,
		policy.uri); err != nil {
		return err
	}

	return nil
}

func (policy CollectionPolicy) Symbol() nft.Symbol {
	return policy.symbol
}

func (policy CollectionPolicy) Name() CollectionName {
	return policy.name
}

func (policy CollectionPolicy) Creator() base.Address {
	return policy.creator
}

func (policy CollectionPolicy) Royalty() nft.PaymentParameter {
	return policy.royalty
}

func (policy CollectionPolicy) Uri() CollectionUri {
	return policy.uri
}

func (policy CollectionPolicy) Equal(cpolicy CollectionPolicy) bool {
	return bytes.Equal(policy.Bytes(), cpolicy.Bytes())
}

func (policy CollectionPolicy) Rebuild() CollectionPolicy {
	return policy
}
