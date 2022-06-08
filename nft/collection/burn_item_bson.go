package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BaseBurnItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"collection": it.collection,
				"nfts":       it.nfts,
				"currency":   it.cid,
			}),
	)
}

type BaseBurnItemBSONUnpacker struct {
	CL string   `bson:"collection"`
	NS bson.Raw `bson:"nfts"`
	CR string   `bson:"currency"`
}

func (it *BaseBurnItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BaseBurnItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.CL, uit.NS, uit.CR)
}
