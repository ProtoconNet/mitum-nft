package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (policy *CollectionPolicy) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	nm string,
	ry uint,
	uri string,
	bws []string,
) error {
	e := util.StringError("failed to unmarshal CollectionPoicy")

	policy.BaseHinter = hint.NewBaseHinter(ht)
	policy.name = CollectionName(nm)
	policy.royalty = PaymentParameter(ry)
	policy.uri = URI(uri)

	whitelist := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e.Wrap(err)
		}
		whitelist[i] = white
	}
	policy.whitelist = whitelist

	return nil
}
