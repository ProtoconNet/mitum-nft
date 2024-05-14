package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it DelegateItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     it.Hint().String(),
			"contract":  it.contract,
			"delegatee": it.delegatee,
			"mode":      it.mode,
			"currency":  it.currency,
		},
	)
}

type DelegateItemBSONUnmarshaler struct {
	Hint      string `bson:"_hint"`
	Contract  string `bson:"contract"`
	Delegatee string `bson:"delegatee"`
	Mode      string `bson:"mode"`
	Currency  string `bson:"currency"`
}

func (it *DelegateItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u DelegateItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unmarshal(enc, ht, u.Contract, u.Delegatee, u.Mode, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
