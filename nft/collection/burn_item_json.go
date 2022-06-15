package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BurnItemJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID           `json:"nft"`
	CR currency.CurrencyID `json:"currency"`
}

func (it BurnItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BurnItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		NF:         it.nft,
		CR:         it.cid,
	})
}

type BurnItemJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	CR string          `json:"currency"`
}

func (it *BurnItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit BurnItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.CR)
}
