package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (nft NFT) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(nft.Hint()),
		bson.M{
			"id":          nft.id,
			"owner":       nft.owner,
			"hash":        nft.hash,
			"uri":         nft.uri,
			"approved":    nft.approved,
			"copyrighter": nft.copyrighter,
		}),
	)
}

type NFTBSONUnpacker struct {
	ID bson.Raw            `bson:"id"`
	ON base.AddressDecoder `bson:"owner"`
	HS string              `bson:"hash"`
	UR string              `bson:"uri"`
	AP base.AddressDecoder `bson:"approved"`
	CP base.AddressDecoder `bson:"copyrighter"`
}

func (nft *NFT) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var unft NFTBSONUnpacker
	if err := enc.Unmarshal(b, &unft); err != nil {
		return err
	}

	return nft.unpack(enc, unft.ID, unft.ON, unft.HS, unft.UR, unft.AP, unft.CP)
}
