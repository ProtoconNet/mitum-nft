package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type MintFormJSONPacker struct {
	jsonenc.HintedHead
	HS nft.NFTHash       `json:"hash"`
	UR nft.URI           `json:"uri"`
	CR []nft.RightHolder `json:"creators"`
	CP []nft.RightHolder `json:"copyrighters"`
}

func (form MintForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(MintFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		HS:         form.hash,
		UR:         form.uri,
		CR:         form.creators,
		CP:         form.copyrighters,
	})
}

type MintFormJSONUnpacker struct {
	HS string          `json:"hash"`
	UR string          `json:"uri"`
	CR json.RawMessage `json:"creators"`
	CP json.RawMessage `json:"copyrighters"`
}

func (form *MintForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufo MintFormJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.HS, ufo.UR, ufo.CR, ufo.CP)
}

type MintItemJSONPacker struct {
	jsonenc.HintedHead
	CL extensioncurrency.ContractID `json:"collection"`
	FO []MintForm                   `json:"forms"`
	CR currency.CurrencyID          `json:"currency"`
}

func (it BaseMintItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(MintItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		CL:         it.collection,
		FO:         it.forms,
		CR:         it.cid,
	})
}

type MintItemJSONUnpacker struct {
	CL string          `json:"collection"`
	FO json.RawMessage `json:"forms"`
	CR string          `json:"currency"`
}

func (it *BaseMintItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit MintItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.CL, uit.FO, uit.CR)
}
