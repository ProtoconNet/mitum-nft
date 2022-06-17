package nft

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (signer Signer) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(signer.Hint()),
		bson.M{
			"account": signer.account,
			"signed":  signer.signed,
		}),
	)
}

type SignerBSONUnpacker struct {
	AC base.AddressDecoder `bson:"account"`
	SG bool                `bson:"signed"`
}

func (signer *Signer) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var us SignerBSONUnpacker
	if err := enc.Unmarshal(b, &us); err != nil {
		return err
	}

	return signer.unpack(enc, us.AC, us.SG)
}
