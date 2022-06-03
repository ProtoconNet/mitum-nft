package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/util/encoder"
)

func (nid *NFTID) unpack(
	enc encoder.Encoder,
	collection string,
	idx uint64,
) error {
	nid.collection = extensioncurrency.ContractID(collection)
	nid.idx = idx

	return nil
}
