package service

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s ServiceDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"design": s.Design,
		},
	)
}

type ServiceDesignStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Design bson.Raw `bson:"design"`
}

func (s *ServiceDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ServiceDesignStateValue")

	var u ServiceDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var sd timestamp.Design
	if err := sd.DecodeBSON(u.Design, enc); err != nil {
		return e(err, "")
	}
	s.Design = sd

	return nil
}

func (s TimeStampLastIndexStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     s.Hint().String(),
			"projectid": s.ProjectID,
			"index":     s.Index,
		},
	)
}

type TimeStampLastIndexStateValueBSONUnmarshaler struct {
	Hint      string `bson:"_hint"`
	ProjectID string `bson:"projectid"`
	Index     uint64 `bson:"index"`
}

func (s *TimeStampLastIndexStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of TimeStampLastIndexStateValue")

	var u TimeStampLastIndexStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	s.ProjectID = u.ProjectID
	s.Index = u.Index

	return nil
}

func (s TimeStampItemStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         s.Hint().String(),
			"timestampitem": s.TimeStampItem,
		},
	)
}

type TimeStampStateValueBSONUnmarshaler struct {
	Hint          string   `bson:"_hint"`
	TimeStampItem bson.Raw `bson:"timestampitem"`
}

func (s *TimeStampItemStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of TimeStampItemStateValue")

	var u TimeStampStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var n timestamp.TimeStampItem
	if err := n.DecodeBSON(u.TimeStampItem, enc); err != nil {
		return e(err, "")
	}
	s.TimeStampItem = n

	return nil
}
