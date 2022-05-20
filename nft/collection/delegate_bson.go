package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact DelegateFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":     fact.h,
				"token":    fact.token,
				"sender":   fact.sender,
				"agents":   fact.agents,
				"mode":     fact.mode,
				"currency": fact.cid,
			}))
}

type DelegateFactBSONUnpacker struct {
	H  valuehash.Bytes       `bson:"hash"`
	TK []byte                `bson:"token"`
	SD base.AddressDecoder   `bson:"sender"`
	AG []base.AddressDecoder `bson:"agents"`
	MD string                `bson:"mode"`
	CR string                `bson:"currency"`
}

func (fact *DelegateFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact DelegateFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.AG, ufact.MD, ufact.CR)
}

func (op *Delegate) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
