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
	Design nft.Design `json:"collectiondesign"`
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

func (s *CollectionStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionDesignStateValue")

	var u CollectionDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var nd nft.Design
	if err := nd.DecodeJSON(u.Design, enc); err != nil {
		return e(err, "")
	}
	s.Design = nd

	return nil
}

type LastNFTIndexStateValueJSONMarshaler struct {
	hint.BaseHinter
	Index nft.NFTID `json:"index"`
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
	Hint  hint.Hint       `json:"_hint"`
	Index json.RawMessage `json:"index"`
}

func (s *LastNFTIndexStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionLastNFTIndexStateValue")

	var u LastNFTIndexStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	hinter, err := enc.Decode(u.Index)
	if err != nil {
		return e(err, "")
	}
	id, ok := hinter.(nft.NFTID)
	if !ok {
		return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
	} else {
		s.id = id
	}

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

type OperatorsBookStateValueJSONMarshaler struct {
	hint.BaseHinter
	Operators OperatorsBook `json:"operatorsbook"`
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

func (s *OperatorsBookStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of OperatorsBookStateValue")

	var u OperatorsBookStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var operators OperatorsBook
	if err := operators.DecodeJSON(u.Operators, enc); err != nil {
		return e(err, "")
	}
	s.Operators = operators

	return nil
}
