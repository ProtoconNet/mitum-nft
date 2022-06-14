package nft

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (r Righter) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(r.Hint()),
		bson.M{
			"account": r.account,
			"signed":  r.signed,
			"clue":    r.clue,
		}),
	)
}

type RighterBSONUnpacker struct {
	AC base.AddressDecoder `bson:"account"`
	SG bool                `bson:"signed"`
	CU string              `bson:"clue"`
}

func (r *Righter) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ur RighterBSONUnpacker
	if err := enc.Unmarshal(b, &ur); err != nil {
		return err
	}

	return r.unpack(enc, ur.AC, ur.SG, ur.CU)
}
