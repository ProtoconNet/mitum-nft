package nft

import (
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type RighterJSONPacker struct {
	jsonenc.HintedHead
	AC base.Address `json:"account"`
	SG bool         `json:"signed"`
	CU string       `json:"clue"`
}

func (r Righter) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(RighterJSONPacker{
		HintedHead: jsonenc.NewHintedHead(r.Hint()),
		AC:         r.account,
		SG:         r.signed,
		CU:         r.clue,
	})
}

type RighterJSONUnpacker struct {
	AC base.AddressDecoder `json:"account"`
	SG bool                `json:"signed"`
	CU string              `json:"clue"`
}

func (r *Righter) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ur RighterJSONUnpacker
	if err := enc.Unmarshal(b, &ur); err != nil {
		return err
	}

	return r.unpack(enc, ur.AC, ur.SG, ur.CU)
}
