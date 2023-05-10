package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (form MintForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        form.Hint().String(),
			"hash":         form.hash,
			"uri":          form.uri,
			"creators":     form.creators,
			"copyrighters": form.copyrighters,
		},
	)
}

type MintFormBSONUnmarshaler struct {
	Hint         string   `bson:"_hint"`
	Hash         string   `bson:"hash"`
	URI          string   `bson:"uri"`
	Creators     bson.Raw `bson:"creators"`
	Copyrighters bson.Raw `bson:"copyrighters"`
}

func (form *MintForm) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of MintForm")

	var u MintFormBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return form.unmarshal(enc, ht, u.Hash, u.URI, u.Creators, u.Copyrighters)
}

func (it MintItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      it.Hint().String(),
			"contract":   it.contract,
			"collection": it.collection,
			"form":       it.form,
			"currency":   it.currency,
		},
	)
}

type MintItemBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Contract   string   `bson:"contract"`
	Collection string   `bson:"collection"`
	Form       bson.Raw `bson:"form"`
	Currency   string   `bson:"currency"`
}

func (it *MintItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of MintItem")

	var u MintItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, ht, u.Contract, u.Collection, u.Form, u.Currency)
}
