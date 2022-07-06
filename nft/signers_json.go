package nft

import (
	"encoding/json"

	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type SignersJSONPacker struct {
	jsonenc.HintedHead
	TT uint     `json:"total"`
	SG []Signer `json:"signers"`
}

func (signers Signers) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SignersJSONPacker{
		HintedHead: jsonenc.NewHintedHead(signers.Hint()),
		TT:         signers.total,
		SG:         signers.signers,
	})
}

type SignersJSONUnpacker struct {
	TT uint            `json:"total"`
	SG json.RawMessage `json:"signers"`
}

func (signers *Signers) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var us SignersJSONUnpacker
	if err := enc.Unmarshal(b, &us); err != nil {
		return err
	}

	return signers.unpack(enc, us.TT, us.SG)
}
