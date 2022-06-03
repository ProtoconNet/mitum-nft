package nft

import (
	"fmt"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxNFTsInCollection = 10000

var (
	NFTIDType   = hint.Type("mitum-nft-nft-id")
	NFTIDHint   = hint.NewHint(NFTIDType, "v0.0.1")
	NFTIDHinter = NFTID{BaseHinter: hint.NewBaseHinter(NFTIDHint)}
)

type NFTID struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	idx        uint64
}

func NewNFTID(collection extensioncurrency.ContractID, idx uint64) NFTID {
	return NFTID{
		BaseHinter: hint.NewBaseHinter(NFTIDHint),
		collection: collection,
		idx:        idx,
	}
}

func MustNewNFTID(collection extensioncurrency.ContractID, idx uint64) NFTID {
	id := NewNFTID(collection, idx)

	if err := id.IsValid(nil); err != nil {
		panic(err)
	}

	return id
}

func (nid NFTID) Bytes() []byte {
	return util.ConcatBytesSlice(
		nid.collection.Bytes(),
		util.Uint64ToBytes(nid.idx),
	)
}

func (nid NFTID) Hint() hint.Hint {
	return NFTIDHint
}

func (nid NFTID) Hash() valuehash.Hash {
	return nid.GenerateHash()
}

func (nid NFTID) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(nid.Bytes())
}

func (nid NFTID) Collection() extensioncurrency.ContractID {
	return nid.collection
}

func (nid NFTID) Idx() uint64 {
	return nid.idx
}

func (nid NFTID) Equal(cnid NFTID) bool {
	return nid.String() == cnid.String()
}

func (nid NFTID) String() string {
	return fmt.Sprintf("%s-%d", nid.collection, nid.idx)
}

func (nid NFTID) IsValid([]byte) error {
	if nid.idx > uint64(MaxNFTsInCollection) {
		return isvalid.InvalidError.Errorf("nid idx over max; %d > %d", nid.idx, MaxNFTsInCollection)
	}
	if err := nid.collection.IsValid(nil); err != nil {
		return err
	}

	return nil
}
