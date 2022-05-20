package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CollectionPolicyJSONPacker struct {
	jsonenc.HintedHead
	SB nft.Symbol           `json:"symbol"`
	NM CollectionName       `json:"name"`
	CE base.Address         `json:"creator"`
	RY nft.PaymentParameter `json:"royalty"`
	UR CollectionUri        `json:"uri"`
}

func (cp CollectionPolicy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CollectionPolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(cp.Hint()),
		SB:         cp.symbol,
		NM:         cp.name,
		CE:         cp.creator,
		RY:         cp.royalty,
		UR:         cp.uri,
	})
}

type CollectionPolicyJSONUnpacker struct {
	SB string              `json:"symbol"`
	NM string              `json:"name"`
	CE base.AddressDecoder `json:"creator"`
	RY uint                `json:"royalty"`
	UR string              `json:"uri"`
}

func (cp *CollectionPolicy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ucp CollectionPolicyJSONUnpacker
	if err := enc.Unmarshal(b, &ucp); err != nil {
		return err
	}

	return cp.unpack(enc, ucp.SB, ucp.NM, ucp.CE, ucp.RY, ucp.UR)
}
