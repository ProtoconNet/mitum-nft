package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (d Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(d.Hint()),
			bson.M{
				"parent":  d.parent,
				"creator": d.creator,
				"symbol":  d.symbol,
				"policy":  d.policy,
			}))
}

type DesignBSONUnpacker struct {
	PR base.AddressDecoder `bson:"parent"`
	CR base.AddressDecoder `bson:"creator"`
	SB string              `bson:"symbol"`
	PO bson.Raw            `bson:"policy"`
}

func (d *Design) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ud DesignBSONUnpacker
	if err := bson.Unmarshal(b, &ud); err != nil {
		return err
	}

	return d.unpack(enc, ud.PR, ud.CR, ud.SB, ud.PO)
}
