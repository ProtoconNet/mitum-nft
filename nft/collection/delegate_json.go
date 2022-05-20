package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type DelegateFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	AG []base.Address      `json:"agents"`
	MD DelegateMode        `json:"mode"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact DelegateFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DelegateFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		AG:         fact.agents,
		MD:         fact.mode,
		CR:         fact.cid,
	})
}

type DelegateFactJSONUnpacker struct {
	H  valuehash.Bytes       `json:"hash"`
	TK []byte                `json:"token"`
	SD base.AddressDecoder   `json:"sender"`
	AG []base.AddressDecoder `json:"agents"`
	MD string                `json:"mode"`
	CR string                `json:"currency"`
}

func (fact *DelegateFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact DelegateFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.AG, ufact.MD, ufact.CR)
}

func (op *Delegate) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
