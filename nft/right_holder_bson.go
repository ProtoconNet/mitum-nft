package nft

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (r RightHolder) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(r.Hint()),
		bson.M{
			"account": r.account,
			"signed":  r.signed,
		}),
	)
}

type RightHolderBSONUnpacker struct {
	AC base.AddressDecoder `bson:"account"`
	SG bool                `bson:"signed"`
}

func (r *RightHolder) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ur RightHolderBSONUnpacker
	if err := enc.Unmarshal(b, &ur); err != nil {
		return err
	}

	return r.unpack(enc, ur.AC, ur.SG)
}
