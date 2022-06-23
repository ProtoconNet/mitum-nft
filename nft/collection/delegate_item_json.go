package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DelegateItemJSONPacker struct {
	jsonenc.HintedHead
	CL extensioncurrency.ContractID `json:"collection"`
	AG base.Address                 `json:"agent"`
	MD DelegateMode                 `json:"mode"`
	CR currency.CurrencyID          `json:"currency"`
}

func (it DelegateItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DelegateItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		CL:         it.collection,
		AG:         it.agent,
		MD:         it.mode,
		CR:         it.cid,
	})
}

type DelegateItemJSONUnpacker struct {
	CL string              `json:"collection"`
	AG base.AddressDecoder `json:"agent"`
	MD string              `json:"mode"`
	CR string              `json:"currency"`
}

func (it *DelegateItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit DelegateItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.CL, uit.AG, uit.MD, uit.CR)
}
