package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CollectionDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CollectionDesign CollectionDesign `json:"collectiondesign"`
}

func (s CollectionDesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		CollectionDesignStateValueJSONMarshaler(s),
	)
}

type CollectionDesignStateValueJSONUnmarshaler struct {
	Hint             hint.Hint       `json:"_hint"`
	CollectionDesign json.RawMessage `json:"collectiondesign"`
}

func (s *CollectionDesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionDesignStateValue")

	var u CollectionDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var cd CollectionDesign
	if err := cd.DecodeJSON(u.CollectionDesign, enc); err != nil {
		return e(err, "")
	}
	s.CollectionDesign = cd

	return nil
}

type CollectionLastNFTIndexStateValueJSONMarshaler struct {
	hint.BaseHinter
	// Collection extensioncurrency.ContractID `json:"collection"`
	Index uint64 `json:"index"`
}

func (s CollectionLastNFTIndexStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		CollectionLastNFTIndexStateValueJSONMarshaler(s),
	)
}

type CollectionLastNFTIndexStateValueJSONUnmarshaler struct {
	Hint hint.Hint `json:"_hint"`
	// Collection string    `json:"collection"`
	Index uint64 `json:"index"`
}

func (s *CollectionLastNFTIndexStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionLastNFTIndexStateValue")

	var u CollectionLastNFTIndexStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)
	// s.Collection = extensioncurrency.ContractID(u.Collection)
	s.Index = u.Index

	return nil
}

type NFTStateValueJSONMarshaler struct {
	hint.BaseHinter
	NFT nft.NFT `json:"nft"`
}

func (s NFTStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		NFTStateValueJSONMarshaler(s),
	)
}

type NFTStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	NFT  json.RawMessage `json:"nft"`
}

func (s *NFTStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTStateValue")

	var u NFTStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var n nft.NFT
	if err := n.DecodeJSON(u.NFT, enc); err != nil {
		return e(err, "")
	}
	s.NFT = n

	return nil
}

type NFTBoxStateValueJSONMarshaler struct {
	hint.BaseHinter
	Box NFTBox `json:"nftbox"`
}

func (s NFTBoxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		NFTBoxStateValueJSONMarshaler(s),
	)
}

type NFTBoxStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	Box  json.RawMessage `json:"nftbox"`
}

func (s *NFTBoxStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTBoxStateValue")

	var u NFTBoxStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var box NFTBox
	if err := box.DecodeJSON(u.Box, enc); err != nil {
		return e(err, "")
	}
	s.Box = box

	return nil
}

type AgentBoxStateValueJSONMarshaler struct {
	hint.BaseHinter
	Box AgentBox `json:"agentbox"`
}

func (s AgentBoxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		AgentBoxStateValueJSONMarshaler(s),
	)
}

type AgentBoxStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	Box  json.RawMessage `json:"agentbox"`
}

func (s *AgentBoxStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AgentBoxStateValue")

	var u AgentBoxStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var box AgentBox
	if err := box.DecodeJSON(u.Box, enc); err != nil {
		return e(err, "")
	}
	s.Box = box

	return nil
}
