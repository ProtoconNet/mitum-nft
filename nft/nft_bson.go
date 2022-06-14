package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (n NFT) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(n.Hint()),
		bson.M{
			"id":           n.id,
			"owner":        n.owner,
			"hash":         n.hash,
			"uri":          n.uri,
			"approved":     n.approved,
			"creators":     n.creators,
			"copyrighters": n.copyrighters,
		}),
	)
}

type NFTBSONUnpacker struct {
	ID bson.Raw            `bson:"id"`
	ON base.AddressDecoder `bson:"owner"`
	HS string              `bson:"hash"`
	UR string              `bson:"uri"`
	AP base.AddressDecoder `bson:"approved"`
	CR bson.Raw            `bson:"creators"`
	CP bson.Raw            `bson:"copyrighters"`
}

func (n *NFT) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var un NFTBSONUnpacker
	if err := enc.Unmarshal(b, &un); err != nil {
		return err
	}

	return n.unpack(enc, un.ID, un.ON, un.HS, un.UR, un.AP, un.CR, un.CP)
}
