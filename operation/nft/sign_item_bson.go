package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it SignItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"nft":      it.nft,
			"currency": it.currency,
		},
	)
}

type SignItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	NFT      uint64 `bson:"nft"`
	Currency string `bson:"currency"`
}

func (it *SignItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u SignItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht, u.Contract, u.NFT, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
