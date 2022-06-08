package collection

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type BurnFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash `json:"hash"`
	TK []byte         `json:"token"`
	SD base.Address   `json:"sender"`
	IT []BurnItem     `json:"items"`
}

func (fact BurnFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BurnFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		IT:         fact.items,
	})
}

type BurnFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	IT json.RawMessage     `json:"items"`
}

func (fact *BurnFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact BurnFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.IT)
}

func (op *Burn) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
