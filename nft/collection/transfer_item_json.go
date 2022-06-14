package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type TransferItemJSONPacker struct {
	jsonenc.HintedHead
	RC base.Address        `json:"receiver"`
	NS []nft.NFTID         `json:"nfts"`
	CR currency.CurrencyID `json:"currency"`
}

func (it BaseTransferItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(TransferItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		RC:         it.receiver,
		NS:         it.nfts,
		CR:         it.cid,
	})
}

type TransferItemJSONUnpacker struct {
	RC base.AddressDecoder `json:"receiver"`
	NS json.RawMessage     `json:"nfts"`
	CR string              `json:"currency"`
}

func (it *BaseTransferItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit TransferItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.RC, uit.NS, uit.CR)
}
