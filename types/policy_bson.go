package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (policy CollectionPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":            policy.Hint().String(),
		"name":             policy.name,
		"royalty":          policy.royalty,
		"uri":              policy.uri,
		"minter_whitelist": policy.whitelist,
	})
}

type PolicyBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Name    string   `bson:"name"`
	Royalty uint     `bson:"royalty"`
	URI     string   `bson:"uri"`
	Whites  []string `bson:"minter_whitelist"`
}

func (policy *CollectionPolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of CollectionPolicy")

	var u PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return policy.unpack(enc, ht, u.Name, u.Royalty, u.URI, u.Whites)
}
