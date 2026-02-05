package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (ob AllApprovedBook) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":        ob.Hint().String(),
		"all_approved": ob.allApproved,
	})
}

type OperatorsBookBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Operators []string `bson:"all_approved"`
}

func (ob *AllApprovedBook) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of all-approved book")

	var u OperatorsBookBSONUnmarshaler
	if err := bsonenc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return ob.unpack(enc, ht, u.Operators)
}
