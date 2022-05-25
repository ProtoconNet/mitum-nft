package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PolicyJSONPacker struct {
	jsonenc.HintedHead
	NM CollectionName       `json:"name"`
	RY nft.PaymentParameter `json:"royalty"`
	UR string               `json:"uri"`
	LI currency.Big         `json:"limit"`
}

func (p Policy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(p.Hint()),
		NM:         p.name,
		RY:         p.royalty,
		UR:         p.uri.String(),
		LI:         p.limit,
	})
}

type PolicyJSONUnpacker struct {
	NM string `json:"name"`
	RY uint   `json:"royalty"`
	UR string `json:"uri"`
	LI string `json:"limit"`
}

func (p *Policy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var up PolicyJSONUnpacker
	if err := enc.Unmarshal(b, &up); err != nil {
		return err
	}

	return p.unpack(enc, up.NM, up.RY, up.UR, up.LI)
}
