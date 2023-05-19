package timestamp

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t TimeStampItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":              t.Hint().String(),
		"projectid":          t.projectID,
		"request_timestamp":  t.requestTimeStamp,
		"response_timestamp": t.responseTimeStamp,
		"timestampid":        t.timestampID,
		"data":               t.data,
	})
}

type TimeStampItemBSONUnmarshaler struct {
	Hint              string `bson:"_hint"`
	ProjectID         string `bson:"projectid"`
	RequestTimeStamp  uint64 `bson:"request_timestamp"`
	ResponseTimeStamp uint64 `bson:"response_timestamp"`
	TimeStampID       uint64 `bson:"timestampid"`
	Data              string `bson:"data"`
}

func (t *TimeStampItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of TimeStampItem")

	var u TimeStampItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return t.unmarshal(enc, ht, u.ProjectID, u.RequestTimeStamp, u.ResponseTimeStamp, u.TimeStampID, u.Data)
}
