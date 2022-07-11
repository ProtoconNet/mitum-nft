package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type CollectionRegisterFormJSONPacker struct {
	jsonenc.HintedHead
	TG base.Address                 `json:"target"`
	SB extensioncurrency.ContractID `json:"symbol"`
	NM CollectionName               `json:"name"`
	RY nft.PaymentParameter         `json:"royalty"`
	UR nft.URI                      `json:"uri"`
	WH []base.Address               `json:"whites"`
}

func (form CollectionRegisterForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CollectionRegisterFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		TG:         form.target,
		SB:         form.symbol,
		NM:         form.name,
		RY:         form.royalty,
		UR:         form.uri,
		WH:         form.whites,
	})
}

type CollectionRegisterFormJSONUnpacker struct {
	TG base.AddressDecoder   `json:"target"`
	SB string                `json:"symbol"`
	NM string                `json:"name"`
	RY uint                  `json:"royalty"`
	UR string                `json:"uri"`
	WH []base.AddressDecoder `json:"whites"`
}

func (form *CollectionRegisterForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uf CollectionRegisterFormJSONUnpacker
	if err := enc.Unmarshal(b, &uf); err != nil {
		return err
	}
	return form.unpack(enc, uf.TG, uf.SB, uf.NM, uf.RY, uf.UR, uf.WH)
}

type CollectionRegisterFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash         `json:"hash"`
	TK []byte                 `json:"token"`
	SD base.Address           `json:"sender"`
	FO CollectionRegisterForm `json:"form"`
	CR currency.CurrencyID    `json:"currency"`
}

func (fact CollectionRegisterFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CollectionRegisterFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		FO:         fact.form,
		CR:         fact.cid,
	})
}

type CollectionRegisterFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	FO json.RawMessage     `json:"form"`
	CR string              `json:"currency"`
}

func (fact *CollectionRegisterFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact CollectionRegisterFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.FO, ufact.CR)
}

func (op *CollectionRegister) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
