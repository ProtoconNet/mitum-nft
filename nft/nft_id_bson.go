package nft

import (
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (nid NFTID) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(nid.Hint()),
		bson.M{
			"collection": nid.collection,
			"idx":        nid.idx,
		}),
	)
}

type NFTIDBSONUnpacker struct {
	CL string `bson:"collection"`
	ID uint64 `bson:"idx"`
}

func (nid *NFTID) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var unid NFTIDBSONUnpacker
	if err := enc.Unmarshal(b, &unid); err != nil {
		return err
	}

	return nid.unpack(enc, unid.CL, unid.ID)
}
