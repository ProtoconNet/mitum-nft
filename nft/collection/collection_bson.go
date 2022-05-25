package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (p Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(p.Hint()),
		bson.M{
			"name":    p.name,
			"royalty": p.royalty,
			"uri":     p.uri.String(),
			"limit":   p.limit.String(),
		},
	))
}

type PolicyBSONUnpacker struct {
	NM string `bson:"name"`
	RY uint   `bson:"royalty"`
	UR string `bson:"uri"`
	LI string `bson:"limit"`
}

func (p Policy) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var up PolicyBSONUnpacker
	if err := enc.Unmarshal(b, &up); err != nil {
		return err
	}

	return p.unpack(enc, up.NM, up.RY, up.UR, up.LI)
}
