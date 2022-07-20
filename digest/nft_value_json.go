package digest

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type NFTValueJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFT     `json:"nft"`
	HT base.Height `json:"height"`
}

func (n NFTValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(n.Hint()),
		NF:         n.nft,
		HT:         n.height,
	})
}

type NFTValueJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	HT base.Height     `json:"height"`
}

func (n *NFTValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva NFTValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	err := n.unpack(enc, uva.NF, uva.HT)
	if err != nil {
		return err
	}
	return nil
}
