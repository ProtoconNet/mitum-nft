package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it ApproveItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"approved": it.approved,
				"nft":      it.nft,
				"currency": it.cid,
			}),
	)
}

type ApproveItemBSONUnpacker struct {
	AP base.AddressDecoder `bson:"approved"`
	NF bson.Raw            `bson:"nft"`
	CR string              `bson:"currency"`
}

func (it *ApproveItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit ApproveItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.AP, uit.NF, uit.CR)
}
