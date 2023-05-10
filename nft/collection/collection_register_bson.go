package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (form CollectionRegisterForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      form.Hint().String(),
			"contract":   form.contract,
			"collection": form.collection,
			"name":       form.name,
			"royalty":    form.royalty,
			"uri":        form.uri,
			"whites":     form.whites,
		})
}

type CollectionRegisterFormBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Contract   string   `bson:"contract"`
	Collection string   `bson:"collection"`
	Name       string   `bson:"name"`
	Royalty    uint     `bson:"royalty"`
	URI        string   `bson:"uri"`
	Whites     []string `bson:"whites"`
}

func (form *CollectionRegisterForm) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionRegisterForm")

	var u CollectionRegisterFormBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return form.unmarshal(enc, ht, u.Contract, u.Collection, u.Name, u.Royalty, u.URI, u.Whites)
}

func (fact CollectionRegisterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":    fact.Hint().String(),
		"hash":     fact.BaseFact.Hash().String(),
		"token":    fact.BaseFact.Token(),
		"sender":   fact.sender,
		"form":     fact.form,
		"currency": fact.currency,
	})
}

type CollectionRegisterFactBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Sender   string   `bson:"sender"`
	Form     bson.Raw `bson:"form"`
	Currency string   `bson:"currency"`
}

func (fact *CollectionRegisterFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionRegisterFact")

	var u currency.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf CollectionRegisterFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unmarshal(enc, uf.Sender, uf.Form, uf.Currency)
}

func (op CollectionRegister) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CollectionRegister) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionRegister")

	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
