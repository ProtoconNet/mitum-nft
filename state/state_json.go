package state

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CollectionDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"collectiondesign"`
}

func (s CollectionStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CollectionDesignStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Design:     s.Design,
	})
}

type CollectionDesignStateValueJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	Design json.RawMessage `json:"collectiondesign"`
}

func (s *CollectionStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CollectionDesignStateValue")

	var u CollectionDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var nd types.Design
	if err := nd.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	s.Design = nd

	return nil
}

type LastNFTIndexStateValueJSONMarshaler struct {
	hint.BaseHinter
	Index uint64 `json:"index"`
}

func (s LastNFTIndexStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		LastNFTIndexStateValueJSONMarshaler{
			BaseHinter: s.BaseHinter,
			Index:      s.id,
		},
	)
}

type LastNFTIndexStateValueJSONUnmarshaler struct {
	Hint  hint.Hint `json:"_hint"`
	Index uint64    `json:"index"`
}

func (s *LastNFTIndexStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CollectionLastNFTIndexStateValue")

	var u LastNFTIndexStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	s.id = u.Index
	return nil
}

type NFTStateValueJSONMarshaler struct {
	hint.BaseHinter
	NFT types.NFT `json:"nft"`
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

func (s *NFTStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of NFTStateValue")

	var u NFTStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var n types.NFT
	if err := n.DecodeJSON(u.NFT, enc); err != nil {
		return e.Wrap(err)
	}
	s.NFT = n

	return nil
}

type OperatorsBookStateValueJSONMarshaler struct {
	hint.BaseHinter
	Operators types.AllApprovedBook `json:"operatorsbook"`
}

func (s OperatorsBookStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		OperatorsBookStateValueJSONMarshaler(s),
	)
}

type OperatorsBookStateValueJSONUnmarshaler struct {
	Hint      hint.Hint       `json:"_hint"`
	Operators json.RawMessage `json:"operatorsbook"`
}

func (s *OperatorsBookStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of OperatorsBookStateValue")

	var u OperatorsBookStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var operators types.AllApprovedBook
	if err := operators.DecodeJSON(u.Operators, enc); err != nil {
		return e.Wrap(err)
	}
	s.Operators = operators

	return nil
}
