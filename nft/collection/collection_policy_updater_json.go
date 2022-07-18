package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type CollectionPolicyUpdaterFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash               `json:"hash"`
	TK []byte                       `json:"token"`
	SD base.Address                 `json:"sender"`
	CL extensioncurrency.ContractID `json:"collection"`
	PO CollectionPolicy             `json:"policy"`
	CR currency.CurrencyID          `json:"currency"`
}

func (fact CollectionPolicyUpdaterFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CollectionPolicyUpdaterFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		CL:         fact.collection,
		PO:         fact.policy,
		CR:         fact.cid,
	})
}

type CollectionPolicyUpdaterFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	CL string              `json:"collection"`
	PO json.RawMessage     `json:"policy"`
	CR string              `json:"currency"`
}

func (fact *CollectionPolicyUpdaterFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact CollectionPolicyUpdaterFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.CL, ufact.PO, ufact.CR)
}

func (op *CollectionPolicyUpdater) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
