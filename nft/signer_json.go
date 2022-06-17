package nft

import (
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type SignerJSONPacker struct {
	jsonenc.HintedHead
	AC base.Address `json:"account"`
	SG bool         `json:"signed"`
}

func (signer Signer) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SignerJSONPacker{
		HintedHead: jsonenc.NewHintedHead(signer.Hint()),
		AC:         signer.account,
		SG:         signer.signed,
	})
}

type SignerJSONUnpacker struct {
	AC base.AddressDecoder `json:"account"`
	SG bool                `json:"signed"`
}

func (signer *Signer) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var us SignerJSONUnpacker
	if err := enc.Unmarshal(b, &us); err != nil {
		return err
	}

	return signer.unpack(enc, us.AC, us.SG)
}
