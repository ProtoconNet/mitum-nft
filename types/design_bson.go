package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    de.Hint().String(),
			"contract": de.contract,
			"creator":  de.creator,
			"active":   de.active,
			"policy":   de.policy,
		})
}

type DesignBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Contract string   `bson:"contract"`
	Creator  string   `bson:"creator"`
	Active   bool     `bson:"active"`
	Policy   bson.Raw `bson:"policy"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Design")

	var u DesignBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return de.unpack(enc, ht, u.Contract, u.Creator, u.Active, u.Policy)
}
