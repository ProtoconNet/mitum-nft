package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BurnItemJSONPacker struct {
	jsonenc.HintedHead
	CL extensioncurrency.ContractID `json:"collection"`
	NS []nft.NFTID                  `json:"nfts"`
	CR currency.CurrencyID          `json:"currency"`
}

func (it BaseBurnItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BurnItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		CL:         it.collection,
		NS:         it.nfts,
		CR:         it.cid,
	})
}

type BurnItemJSONUnpacker struct {
	CL string          `json:"collection"`
	NS json.RawMessage `json:"nfts"`
	CR string          `json:"currency"`
}

func (it *BaseBurnItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit BurnItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.CL, uit.NS, uit.CR)
}
