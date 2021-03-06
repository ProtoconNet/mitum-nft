package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CollectionPolicyJSONPacker struct {
	jsonenc.HintedHead
	NM CollectionName       `json:"name"`
	RY nft.PaymentParameter `json:"royalty"`
	UR nft.URI              `json:"uri"`
	WH []base.Address       `json:"whites"`
}

func (p CollectionPolicy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CollectionPolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(p.Hint()),
		NM:         p.name,
		RY:         p.royalty,
		UR:         p.uri,
		WH:         p.whites,
	})
}

type CollectionPolicyJSONUnpacker struct {
	NM string                `json:"name"`
	RY uint                  `json:"royalty"`
	UR string                `json:"uri"`
	WH []base.AddressDecoder `json:"whites"`
}

func (p *CollectionPolicy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var up CollectionPolicyJSONUnpacker
	if err := enc.Unmarshal(b, &up); err != nil {
		return err
	}

	return p.unpack(enc, up.NM, up.RY, up.UR, up.WH)
}
