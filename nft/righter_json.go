package nft

import (
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type RightHolerJSONPacker struct {
	jsonenc.HintedHead
	AC base.Address `json:"account"`
	SG bool         `json:"signed"`
	CU string       `json:"clue"`
}

func (r RightHoler) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(RightHolerJSONPacker{
		HintedHead: jsonenc.NewHintedHead(r.Hint()),
		AC:         r.account,
		SG:         r.signed,
		CU:         r.clue,
	})
}

type RightHolerJSONUnpacker struct {
	AC base.AddressDecoder `json:"account"`
	SG bool                `json:"signed"`
	CU string              `json:"clue"`
}

func (r *RightHoler) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ur RightHolerJSONUnpacker
	if err := enc.Unmarshal(b, &ur); err != nil {
		return err
	}

	return r.unpack(enc, ur.AC, ur.SG, ur.CU)
}
