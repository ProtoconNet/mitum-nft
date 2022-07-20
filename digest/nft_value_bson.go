package digest

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (n NFTValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(n.Hint()),
		bson.M{
			"nft":    n.nft,
			"height": n.height,
		},
	))
}

type NFTValueBSONUnpacker struct {
	NF bson.Raw    `bson:"nft"`
	HT base.Height `bson:"height"`
}

func (n *NFTValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uva NFTValueBSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	return n.unpack(enc, uva.NF, uva.HT)
}
