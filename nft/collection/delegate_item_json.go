package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DelegateItemJSONPacker struct {
	jsonenc.HintedHead
	AG base.Address        `json:"agent"`
	MD DelegateMode        `json:"mode"`
	CR currency.CurrencyID `json:"currency"`
}

func (it DelegateItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DelegateItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		AG:         it.agent,
		MD:         it.mode,
		CR:         it.cid,
	})
}

type DelegateItemJSONUnpacker struct {
	AG base.AddressDecoder `json:"agent"`
	MD string              `json:"mode"`
	CR string              `json:"currency"`
}

func (it *DelegateItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit DelegateItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.AG, uit.MD, uit.CR)
}
