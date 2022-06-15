package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type ApproveItemJSONPacker struct {
	jsonenc.HintedHead
	AP base.Address        `json:"approved"`
	NF nft.NFTID           `json:"nft"`
	CR currency.CurrencyID `json:"currency"`
}

func (it ApproveItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(ApproveItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		AP:         it.approved,
		NF:         it.nft,
		CR:         it.cid,
	})
}

type ApproveItemJSONUnpacker struct {
	AP base.AddressDecoder `json:"approved"`
	NF json.RawMessage     `json:"nft"`
	CR string              `json:"currency"`
}

func (it *ApproveItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit ApproveItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.AP, uit.NF, uit.CR)
}
