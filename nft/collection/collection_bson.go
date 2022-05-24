package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (p Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(p.Hint()),
		bson.M{
			"name":    p.name,
			"creator": p.creator,
			"royalty": p.royalty,
			"uri":     p.uri.String(),
		},
	))
}

type CollectionPolicyBSONUnpacker struct {
	NM string              `bson:"name"`
	CE base.AddressDecoder `bson:"creator"`
	RY uint                `bson:"royalty"`
	UR string              `bson:"uri"`
}

func (p Policy) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ucp CollectionPolicyBSONUnpacker
	if err := enc.Unmarshal(b, &ucp); err != nil {
		return err
	}

	return p.unpack(enc, ucp.NM, ucp.CE, ucp.RY, ucp.UR)
}
