package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type NFTIDJSONPacker struct {
	jsonenc.HintedHead
	CL extensioncurrency.ContractID `json:"collection"`
	ID uint64                       `json:"idx"`
}

func (nid NFTID) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTIDJSONPacker{
		HintedHead: jsonenc.NewHintedHead(nid.Hint()),
		CL:         nid.collection,
		ID:         nid.idx,
	})
}

type NFTIDJSONUnpacker struct {
	CL string `json:"collection"`
	ID uint64 `json:"idx"`
}

func (nid *NFTID) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var unid NFTIDJSONUnpacker
	if err := enc.Unmarshal(b, &unid); err != nil {
		return err
	}

	return nid.unpack(enc, unid.CL, unid.ID)
}
