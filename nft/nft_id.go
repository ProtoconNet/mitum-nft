package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum2/util"
)

var MaxNFTIndex uint64 = 10000

type NFTID uint64

func (nid NFTID) IsValid([]byte) error {
	if nid.Index() > MaxNFTIndex {
		return util.ErrInvalid.Errorf("nft-id index over max, %d > %d", nid.Index(), MaxNFTIndex)
	}

	if nid.Index() == 0 {
		return util.ErrInvalid.Errorf("zero nft-id index, %q", nid)
	}

	// if err := nid.collection.IsValid(nil); err != nil {
	// 	return err
	// }

	return nil
}

func (nid NFTID) Bytes() []byte {
	return util.ConcatBytesSlice(
		util.Uint64ToBytes(nid.Index()),
	)
}

func (nid NFTID) Index() uint64 {
	return uint64(nid)
}

func (nid NFTID) Equal(id NFTID) bool {
	return nid.String() == id.String()
}

func (nid NFTID) String() string {
	index := strconv.FormatUint(nid.Index(), 10)

	l := len(strconv.FormatUint(uint64(MaxNFTIndex), 10)) - len(index)
	for i := 0; i < l; i++ {
		index = "0" + index
	}

	return index
}
