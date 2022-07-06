package nft

import (
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (signers Signers) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(signers.Hint()),
		bson.M{
			"total":   signers.total,
			"signers": signers.signers,
		}),
	)
}

type SignersBSONUnpacker struct {
	TT uint     `bson:"total"`
	SG bson.Raw `bson:"signers"`
}

func (signers *Signers) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var us SignersBSONUnpacker
	if err := enc.Unmarshal(b, &us); err != nil {
		return err
	}

	return signers.unpack(enc, us.TT, us.SG)
}
