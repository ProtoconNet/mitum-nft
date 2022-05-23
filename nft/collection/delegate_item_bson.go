package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it DelegateItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"agent":    it.agent,
				"mode":     it.mode,
				"currency": it.cid,
			}),
	)
}

type DelegateItemBSONUnpacker struct {
	AG base.AddressDecoder `bson:"agent"`
	MD string              `bson:"mode"`
	CR string              `bson:"currency"`
}

func (it *DelegateItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit DelegateItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.AG, uit.MD, uit.CR)
}
