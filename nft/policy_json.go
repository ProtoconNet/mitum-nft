package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DesignJSONPacker struct {
	jsonenc.HintedHead
	SB extension.ContractID `json:"symbol"`
	PO BasePolicy           `json:"policy"`
}

func (d Design) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DesignJSONPacker{
		HintedHead: jsonenc.NewHintedHead(d.Hint()),
		SB:         d.symbol,
		PO:         d.policy,
	})
}

type DesignJSONUnpacker struct {
	SB string          `json:"symbol"`
	PO json.RawMessage `json:"policy"`
}

func (d *Design) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ud DesignJSONUnpacker
	if err := enc.Unmarshal(b, &ud); err != nil {
		return err
	}

	return d.unpack(enc, ud.SB, ud.PO)
}
