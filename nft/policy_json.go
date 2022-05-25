package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DesignJSONPacker struct {
	jsonenc.HintedHead
	PR base.Address         `json:"parent"`
	CR base.Address         `json:"creator"`
	SB extension.ContractID `json:"symbol"`
	PO BasePolicy           `json:"policy"`
}

func (d Design) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DesignJSONPacker{
		HintedHead: jsonenc.NewHintedHead(d.Hint()),
		PR:         d.parent,
		CR:         d.creator,
		SB:         d.symbol,
		PO:         d.policy,
	})
}

type DesignJSONUnpacker struct {
	PR base.AddressDecoder `json:"parent"`
	CR base.AddressDecoder `json:"creator"`
	SB string              `json:"symbol"`
	PO json.RawMessage     `json:"policy"`
}

func (d *Design) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ud DesignJSONUnpacker
	if err := enc.Unmarshal(b, &ud); err != nil {
		return err
	}

	return d.unpack(enc, ud.PR, ud.CR, ud.SB, ud.PO)
}
