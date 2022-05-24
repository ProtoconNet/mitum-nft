package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (d Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(d.Hint()),
			bson.M{
				"symbol": d.symbol,
				"policy": d.policy,
			}))
}

type DesignBSONUnpacker struct {
	SB string   `bson:"symbol"`
	PO bson.Raw `bson:"policy"`
}

func (d *Design) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ud DesignBSONUnpacker
	if err := bson.Unmarshal(b, &ud); err != nil {
		return err
	}

	return d.unpack(enc, ud.SB, ud.PO)
}
