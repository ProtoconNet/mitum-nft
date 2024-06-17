package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it TransferItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"receiver": it.receiver,
			"nft_idx":  it.nftIdx,
			"currency": it.currency,
		},
	)
}

type TransferItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	Receiver string `bson:"receiver"`
	NFTIdx   uint64 `bson:"nft_idx"`
	Currency string `bson:"currency"`
}

func (it *TransferItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u TransferItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht, u.Contract, u.Receiver, u.NFTIdx, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
