package nft

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type NFTJSONPacker struct {
	jsonenc.HintedHead
	ID NFTID        `json:"id"`
	AC bool         `json:"active"`
	ON base.Address `json:"owner"`
	HS NFTHash      `json:"hash"`
	UR URI          `json:"uri"`
	AP base.Address `json:"approved"`
	CR Signers      `json:"creators"`
	CP Signers      `json:"copyrighters"`
}

func (n NFT) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTJSONPacker{
		HintedHead: jsonenc.NewHintedHead(n.Hint()),
		ID:         n.id,
		AC:         n.active,
		ON:         n.owner,
		HS:         n.hash,
		UR:         n.uri,
		AP:         n.approved,
		CR:         n.creators,
		CP:         n.copyrighters,
	})
}

type NFTJSONUnpacker struct {
	ID json.RawMessage     `json:"id"`
	AC bool                `json:"active"`
	ON base.AddressDecoder `json:"owner"`
	HS string              `json:"hash"`
	UR string              `json:"uri"`
	AP base.AddressDecoder `json:"approved"`
	CR json.RawMessage     `json:"creators"`
	CP json.RawMessage     `json:"copyrighters"`
}

func (n *NFT) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var un NFTJSONUnpacker
	if err := enc.Unmarshal(b, &un); err != nil {
		return err
	}

	return n.unpack(enc, un.ID, un.AC, un.ON, un.HS, un.UR, un.AP, un.CR, un.CP)
}
