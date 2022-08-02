package nft

import (
	"fmt"
	"strconv"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxNFTIdx = 10000

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

func NewNFTID(symbol extensioncurrency.ContractID, idx uint64) NFTID {
	return NFTID{
		BaseHinter: hint.NewBaseHinter(NFTIDHint),
		collection: symbol,
		idx:        idx,
	}
}

func MustNewNFTID(symbol extensioncurrency.ContractID, idx uint64) NFTID {
	id := NewNFTID(symbol, idx)

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
	idx := strconv.FormatUint(nid.idx, 10)

	l := len(strconv.FormatUint(uint64(MaxNFTIdx), 10)) - len(idx)
	for i := 0; i < l; i++ {
		idx = "0" + idx
	}

	return fmt.Sprintf("%s-%s", nid.collection, idx)
}

func (nid NFTID) IsValid([]byte) error {
	if nid.idx > uint64(MaxNFTIdx) {
		return isvalid.InvalidError.Errorf("nid idx over max; %d > %d", nid.idx, MaxNFTIdx)
	}

	if nid.idx == 0 {
		return isvalid.InvalidError.Errorf("nid idx must be over zero; %q", nid)
	}

	if err := nid.collection.IsValid(nil); err != nil {
		return err
	}

	return nil
}
