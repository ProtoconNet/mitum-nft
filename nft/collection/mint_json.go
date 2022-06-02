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

type MintFormJSONPacker struct {
	jsonenc.HintedHead
	HS nft.NFTHash  `json:"hash"`
	UR string       `json:"uri"`
	CP base.Address `json:"copyrighter"`
}

func (form MintForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(MintFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		HS:         form.hash,
		UR:         form.uri.String(),
		CP:         form.copyrighter,
	})
}

type MintFormJSONUnpacker struct {
	HS string              `json:"hash"`
	UR string              `json:"uri"`
	CP base.AddressDecoder `json:"copyrighter"`
}

func (form *MintForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufo MintFormJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.HS, ufo.UR, ufo.CP)
}

type MintFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash               `json:"hash"`
	TK []byte                       `json:"token"`
	SD base.Address                 `json:"sender"`
	CL extensioncurrency.ContractID `json:"collection"`
	FO MintForm                     `json:"form"`
	CR currency.CurrencyID          `json:"currency"`
}

func (fact MintFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(MintFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		CL:         fact.collection,
		FO:         fact.form,
		CR:         fact.cid,
	})
}

type MintFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	CL string              `json:"collection"`
	FO json.RawMessage     `json:"form"`
	CR string              `json:"currency"`
}

func (fact *MintFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact MintFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.CL, ufact.FO, ufact.CR)
}

func (op *Mint) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
