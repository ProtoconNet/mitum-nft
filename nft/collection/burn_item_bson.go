package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BurnItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"nft":      it.nft,
				"currency": it.cid,
			}),
	)
}

type BurnItemBSONUnpacker struct {
	NF bson.Raw `bson:"nft"`
	CR string   `bson:"currency"`
}

func (it *BurnItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BurnItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.CR)
}
