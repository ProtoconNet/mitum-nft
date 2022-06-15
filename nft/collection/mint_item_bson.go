package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (form MintForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"hash":         form.hash,
				"uri":          form.uri,
				"creators":     form.creators,
				"copyrighters": form.copyrighters,
			}))
}

type MintFormBSONUnpacker struct {
	HS string   `bson:"hash"`
	UR string   `bson:"uri"`
	CR bson.Raw `bson:"creators"`
	CP bson.Raw `bson:"copyrighters"`
}

func (form *MintForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufo MintFormBSONUnpacker
	if err := bson.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.HS, ufo.UR, ufo.CR, ufo.CP)
}

func (it MintItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"collection": it.collection,
				"form":       it.form,
				"currency":   it.cid,
			}),
	)
}

type MintItemBSONUnpacker struct {
	CL string   `bson:"collection"`
	FO bson.Raw `bson:"form"`
	CR string   `bson:"currency"`
}

func (it *MintItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit MintItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.CL, uit.FO, uit.CR)
}
