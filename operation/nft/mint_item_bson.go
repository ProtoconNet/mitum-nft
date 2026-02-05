package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (it MintItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"receiver": it.receiver,
			"hash":     it.hash,
			"uri":      it.uri,
			"creators": it.creators,
			"currency": it.currency,
		},
	)
}

type MintItemBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Contract string   `bson:"contract"`
	Receiver string   `bson:"receiver"`
	Hash     string   `bson:"hash"`
	Uri      string   `bson:"uri"`
	Creators bson.Raw `bson:"creators"`
	Currency string   `bson:"currency"`
}

func (it *MintItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u MintItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht, u.Contract, u.Receiver, u.Hash, u.Uri, u.Creators, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
