package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type SignItemJSONPacker struct {
	jsonenc.HintedHead
	QU Qualification       `json:"qualification"`
	NF nft.NFTID           `json:"nft"`
	CR currency.CurrencyID `json:"currency"`
}

func (it SignItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SignItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		QU:         it.qualification,
		NF:         it.nft,
		CR:         it.cid,
	})
}

type SignItemJSONUnpacker struct {
	QU string          `json:"qualification"`
	NF json.RawMessage `json:"nft"`
	CR string          `json:"currency"`
}

func (it *SignItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit SignItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.QU, uit.NF, uit.CR)
}
