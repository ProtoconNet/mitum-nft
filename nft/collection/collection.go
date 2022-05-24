package collection

import (
	"net/url"
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
	PolicyType   = hint.Type("mitum-nft-policy")
	PolicyHint   = hint.NewHint(PolicyType, "v0.0.1")
	PolicyHinter = Policy{BaseHinter: hint.NewBaseHinter(PolicyHint)}
)

type Policy struct {
	hint.BaseHinter
	name    CollectionName
	creator base.Address
	royalty nft.PaymentParameter
	uri     url.URL
}

func NewPolicy(name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri url.URL) Policy {
	return Policy{
		BaseHinter: hint.NewBaseHinter(PolicyHint),
		name:       name,
		creator:    creator,
		royalty:    royalty,
		uri:        uri,
	}
}

func MustNewPolicy(name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri url.URL) Policy {
	policy := NewPolicy(name, creator, royalty, uri)

	if err := policy.IsValid(nil); err != nil {
		panic(err)
	}

	return policy
}

func (policy Policy) Bytes() []byte {
	return util.ConcatBytesSlice(
		policy.name.Bytes(),
		policy.creator.Bytes(),
		policy.royalty.Bytes(),
		[]byte(policy.uri.String()),
	)
}

func (policy Policy) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		policy.name,
		policy.creator,
		policy.royalty); err != nil {
		return err
	}

	return nil
}

func (policy Policy) Name() CollectionName {
	return policy.name
}

func (policy Policy) Creator() base.Address {
	return policy.creator
}

func (policy Policy) Royalty() nft.PaymentParameter {
	return policy.royalty
}

func (policy Policy) Uri() url.URL {
	return policy.uri
}

func (policy Policy) Rebuild() nft.BasePolicy {
	return policy
}
