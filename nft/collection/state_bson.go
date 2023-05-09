package collection

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s CollectionDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            s.Hint().String(),
			"collectiondesign": s.CollectionDesign,
		},
	)
}

type CollectionDesignStateValueBSONUnmarshaler struct {
	Hint             string   `bson:"_hint"`
	CollectionDesign bson.Raw `bson:"collectiondesign"`
}

func (s *CollectionDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionDesignStateValue")

	var u CollectionDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var cd CollectionDesign
	if err := cd.DecodeBSON(u.CollectionDesign, enc); err != nil {
		return e(err, "")
	}
	s.CollectionDesign = cd

	return nil
}

func (s CollectionLastNFTIndexStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": s.Hint().String(),
			// "collection": s.Collection,
			"index": s.Index,
		},
	)
}

type CollectionLastNFTIndexStateValueBSONUnmarshaler struct {
	Hint string `bson:"_hint"`
	// Collection string `bson:"collection"`
	Index uint64 `bson:"index"`
}

func (s *CollectionLastNFTIndexStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionLastNFTIndexStateValue")

	var u CollectionLastNFTIndexStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	// s.Collection = extensioncurrency.ContractID(u.Collection)
	s.Index = u.Index

	return nil
}

func (s NFTStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": s.Hint().String(),
			"nft":   s.NFT,
		},
	)
}

type NFTStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	NFT  bson.Raw `bson:"nft"`
}

func (s *NFTStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of NFTStateValue")

	var u NFTStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var n nft.NFT
	if err := n.DecodeBSON(u.NFT, enc); err != nil {
		return e(err, "")
	}
	s.NFT = n

	return nil
}

func (s NFTBoxStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"nftbox": s.Box,
		},
	)
}

type NFTBoxStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	Box  bson.Raw `bson:"nftbox"`
}

func (s *NFTBoxStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of NFTBoxStateValue")

	var u NFTBoxStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var box NFTBox
	if err := box.DecodeBSON(u.Box, enc); err != nil {
		return e(err, "")
	}
	s.Box = box

	return nil
}

func (s AgentBoxStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    s.Hint().String(),
			"agentbox": s.Box,
		},
	)
}

type AgentBoxStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	Box  bson.Raw `bson:"agentbox"`
}

func (s *AgentBoxStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AgentBoxStateValue")

	var u AgentBoxStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var box AgentBox
	if err := box.DecodeBSON(u.Box, enc); err != nil {
		return e(err, "")
	}
	s.Box = box

	return nil
}
