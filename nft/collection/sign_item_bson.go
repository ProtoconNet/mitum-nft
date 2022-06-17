package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it SignItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"qualification": it.qualification,
				"nft":           it.nft,
				"currency":      it.cid,
			}),
	)
}

type SignItemBSONUnpacker struct {
	QU string   `bson:"qualification"`
	NF bson.Raw `bson:"nft"`
	CR string   `bson:"currency"`
}

func (it *SignItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit SignItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.QU, uit.NF, uit.CR)
}
