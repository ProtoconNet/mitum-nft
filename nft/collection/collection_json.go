package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PolicyJSONPacker struct {
	jsonenc.HintedHead
	NM CollectionName       `json:"name"`
	CE base.Address         `json:"creator"`
	RY nft.PaymentParameter `json:"royalty"`
	UR string               `json:"uri"`
}

func (p Policy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(p.Hint()),
		NM:         p.name,
		CE:         p.creator,
		RY:         p.royalty,
		UR:         p.uri.String(),
	})
}

type PolicyJSONUnpacker struct {
	NM string              `json:"name"`
	CE base.AddressDecoder `json:"creator"`
	RY uint                `json:"royalty"`
	UR string              `json:"uri"`
}

func (p *Policy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var up PolicyJSONUnpacker
	if err := enc.Unmarshal(b, &up); err != nil {
		return err
	}

	return p.unpack(enc, up.NM, up.CE, up.RY, up.UR)
}
