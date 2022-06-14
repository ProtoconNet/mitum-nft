package nft

import (
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type RightHolderJSONPacker struct {
	jsonenc.HintedHead
	AC base.Address `json:"account"`
	SG bool         `json:"signed"`
}

func (r RightHolder) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(RightHolderJSONPacker{
		HintedHead: jsonenc.NewHintedHead(r.Hint()),
		AC:         r.account,
		SG:         r.signed,
	})
}

type RightHolderJSONUnpacker struct {
	AC base.AddressDecoder `json:"account"`
	SG bool                `json:"signed"`
}

func (r *RightHolder) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ur RightHolderJSONUnpacker
	if err := enc.Unmarshal(b, &ur); err != nil {
		return err
	}

	return r.unpack(enc, ur.AC, ur.SG)
}
