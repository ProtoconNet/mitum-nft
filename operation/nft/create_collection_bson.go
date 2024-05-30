package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CreateCollectionFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":     fact.Hint().String(),
		"hash":      fact.BaseFact.Hash().String(),
		"token":     fact.BaseFact.Token(),
		"sender":    fact.sender,
		"contract":  fact.contract,
		"name":      fact.name,
		"royalty":   fact.royalty,
		"uri":       fact.uri,
		"whitelist": fact.whitelist,
		"currency":  fact.currency,
	})
}

type CreateCollectionFactBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Sender    string   `bson:"sender"`
	Contract  string   `bson:"contract"`
	Name      string   `bson:"name"`
	Royalty   uint     `bson:"royalty"`
	URI       string   `bson:"uri"`
	Whitelist []string `bson:"whitelist"`
	Currency  string   `bson:"currency"`
}

func (fact *CreateCollectionFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u common.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf CreateCollectionFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	if err := fact.unmarshal(enc, uf.Sender, uf.Contract, uf.Name, uf.Royalty, uf.URI, uf.Whitelist, uf.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op CreateCollection) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CreateCollection) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
