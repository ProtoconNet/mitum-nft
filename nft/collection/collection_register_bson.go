package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form CollectionRegisterForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"target":  form.target,
				"symbol":  form.symbol,
				"name":    form.name,
				"royalty": form.royalty,
				"uri":     form.uri,
				"whites":  form.whites,
			}))
}

type CollectionRegisterFormBSONUnpacker struct {
	TG base.AddressDecoder   `bson:"target"`
	SB string                `bson:"symbol"`
	NM string                `bson:"name"`
	RY uint                  `bson:"royalty"`
	UR string                `bson:"uri"`
	WH []base.AddressDecoder `bson:"whites"`
}

func (form *CollectionRegisterForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uf CollectionRegisterFormBSONUnpacker
	if err := bson.Unmarshal(b, &uf); err != nil {
		return err
	}

	return form.unpack(enc, uf.TG, uf.SB, uf.NM, uf.RY, uf.UR, uf.WH)
}

func (fact CollectionRegisterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":     fact.h,
				"token":    fact.token,
				"sender":   fact.sender,
				"form":     fact.form,
				"currency": fact.cid,
			}))
}

type CollectionRegisterFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	FO bson.Raw            `bson:"form"`
	CR string              `bson:"currency"`
}

func (fact *CollectionRegisterFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact CollectionRegisterFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.FO, ufact.CR)
}

func (op *CollectionRegister) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
