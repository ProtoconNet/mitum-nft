package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type MintFormJSONPacker struct {
	jsonenc.HintedHead
	HS nft.NFTHash  `json:"hash"`
	UR nft.URI      `json:"uri"`
	CP base.Address `json:"copyrighter"`
}

func (form MintForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(MintFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		HS:         form.hash,
		UR:         form.uri,
		CP:         form.copyrighter,
	})
}

type MintFormJSONUnpacker struct {
	HS string `json:"hash"`
	UR string `json:"uri"`
	CP string `json:"copyrighter"`
}

func (form *MintForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufo MintFormJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.HS, ufo.UR, ufo.CP)
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
	var utn MintItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &utn); err != nil {
		return err
	}

	return it.unpack(enc, utn.CL, utn.FO, utn.CR)
}
