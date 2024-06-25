package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CollectionPolicyJSONMarshaler struct {
	hint.BaseHinter
	Name      CollectionName   `json:"name"`
	Royalty   PaymentParameter `json:"royalty"`
	URI       URI              `json:"uri"`
	Whitelist []base.Address   `json:"minter_whitelist"`
}

func (policy CollectionPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionPolicyJSONMarshaler{
		BaseHinter: policy.BaseHinter,
		Name:       policy.name,
		Royalty:    policy.royalty,
		URI:        policy.uri,
		Whitelist:  policy.whitelist,
	})
}

type CollectionPolicyJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Name      string    `json:"name"`
	Royalty   uint      `json:"royalty"`
	URI       string    `json:"uri"`
	Whitelist []string  `json:"minter_whitelist"`
}

func (policy *CollectionPolicy) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CollectionPolicy")

	var u CollectionPolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return policy.unpack(enc, u.Hint, u.Name, u.Royalty, u.URI, u.Whitelist)
}
