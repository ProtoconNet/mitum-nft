package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type ApproveFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	AP base.Address        `json:"approved"`
	NS []nft.NFTID         `json:"nfts"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact ApproveFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(ApproveFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		AP:         fact.approved,
		NS:         fact.nfts,
		CR:         fact.cid,
	})
}

type ApproveFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	AP base.AddressDecoder `json:"approved"`
	NS json.RawMessage     `json:"nfts"`
	CR string              `json:"currency"`
}

func (fact *ApproveFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact ApproveFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.AP, ufact.NS, ufact.CR)
}

func (op *Approve) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
