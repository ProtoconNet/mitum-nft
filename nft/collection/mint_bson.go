package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form MintForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"hash":        form.hash,
				"uri":         form.uri.String(),
				"copyrighter": form.copyrighter,
			}))
}

type MintFormBSONUnpacker struct {
	HS string              `bson:"hash"`
	UR string              `bson:"uri"`
	CP base.AddressDecoder `bson:"copyrighter"`
}

func (form *MintForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufo MintFormBSONUnpacker
	if err := bson.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.HS, ufo.UR, ufo.CP)
}

func (fact MintFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":       fact.h,
				"token":      fact.token,
				"sender":     fact.sender,
				"collection": fact.collection,
				"form":       fact.form,
				"currency":   fact.cid,
			}))
}

type MintFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	CL string              `bson:"collection"`
	FO bson.Raw            `bson:"form"`
	CR string              `bson:"currency"`
}

func (fact *MintFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact MintFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.CL, ufact.FO, ufact.CR)
}

func (op *Mint) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
