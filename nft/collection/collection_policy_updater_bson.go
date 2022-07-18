package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact CollectionPolicyUpdaterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":       fact.h,
				"token":      fact.token,
				"sender":     fact.sender,
				"collection": fact.collection,
				"policy":     fact.policy,
				"currency":   fact.cid,
			}))
}

type CollectionPolicyUpdaterFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	CL string              `bson:"collection"`
	PO bson.Raw            `bson:"policy"`
	CR string              `bson:"currency"`
}

func (fact *CollectionPolicyUpdaterFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact CollectionPolicyUpdaterFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.CL, ufact.PO, ufact.CR)
}

func (op *CollectionPolicyUpdater) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
