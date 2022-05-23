package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BaseApproveItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"approved": it.approved,
				"nfts":     it.nfts,
				"currency": it.cid,
			}),
	)
}

type BaseApproveItemBSONUnpacker struct {
	AP base.AddressDecoder `bson:"approved"`
	NS bson.Raw            `bson:"nfts"`
	CR string              `bson:"currency"`
}

func (it *BaseApproveItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BaseApproveItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.AP, uit.NS, uit.CR)
}
