package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it TransferItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"receiver": it.receiver,
				"nft":      it.nft,
				"currency": it.cid,
			}),
	)
}

type TransferItemBSONUnpacker struct {
	RC base.AddressDecoder `bson:"receiver"`
	NF bson.Raw            `bson:"nft"`
	CR string              `bson:"currency"`
}

func (it *TransferItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit TransferItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.RC, uit.NF, uit.CR)
}
