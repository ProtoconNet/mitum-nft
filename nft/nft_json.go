package nft

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type NFTJSONPacker struct {
	jsonenc.HintedHead
	ID NFTID        `json:"id"`
	ON base.Address `json:"owner"`
	HS NFTHash      `json:"hash"`
	UR URI          `json:"uri"`
	AP base.Address `json:"approved"`
	CP base.Address `json:"copyrighter"`
}

func (nft NFT) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTJSONPacker{
		HintedHead: jsonenc.NewHintedHead(nft.Hint()),
		ID:         nft.id,
		ON:         nft.owner,
		HS:         nft.hash,
		UR:         nft.uri,
		AP:         nft.approved,
		CP:         nft.copyrighter,
	})
}

type NFTJSONUnpacker struct {
	ID json.RawMessage     `json:"id"`
	ON base.AddressDecoder `json:"owner"`
	HS string              `json:"hash"`
	UR string              `json:"uri"`
	AP base.AddressDecoder `json:"approved"`
	CP base.AddressDecoder `json:"copyrighter"`
}

func (nft *NFT) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var unft NFTJSONUnpacker
	if err := enc.Unmarshal(b, &unft); err != nil {
		return err
	}

	return nft.unpack(enc, unft.ID, unft.ON, unft.HS, unft.UR, unft.AP, unft.CP)
}
