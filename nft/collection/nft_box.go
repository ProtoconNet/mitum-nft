package collection

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var NFTBoxHint = hint.MustNewHint("mitum-nft-nft-box-v0.0.1")

type NFTBox struct {
	hint.BaseHinter
	nfts []nft.NFTID
}

func NewNFTBox(nfts []nft.NFTID) NFTBox {
	ns := []nft.NFTID{}

	if nfts != nil {
		ns = nfts
	}

	return NFTBox{BaseHinter: hint.NewBaseHinter(NFTBoxHint), nfts: ns}
}

func (nbx NFTBox) Bytes() []byte {
	bns := make([][]byte, len(nbx.nfts))
	for i, n := range nbx.nfts {
		bns[i] = n.Bytes()
	}

	return util.ConcatBytesSlice(bns...)
}

func (nbx NFTBox) Hint() hint.Hint {
	return NFTBoxHint
}

func (nbx NFTBox) Hash() util.Hash {
	return nbx.GenerateHash()
}

func (nbx NFTBox) GenerateHash() util.Hash {
	return valuehash.NewSHA256(nbx.Bytes())
}

func (nbx NFTBox) IsEmpty() bool {
	return len(nbx.nfts) < 1
}

func (nbx NFTBox) IsValid([]byte) error {
	for _, n := range nbx.nfts {
		if err := n.IsValid(nil); err != nil {
			return err
		}
	}
	return nil
}

func (nbx NFTBox) Equal(b NFTBox) bool {
	nbx.Sort(true)
	b.Sort(true)
	for i := range nbx.nfts {
		if !nbx.nfts[i].Equal(b.nfts[i]) {
			return false
		}
	}
	return true
}

func (nbx *NFTBox) Sort(ascending bool) {
	sort.Slice(nbx.nfts, func(i, j int) bool {
		if ascending {
			return bytes.Compare(nbx.nfts[j].Bytes(), nbx.nfts[i].Bytes()) > 0
		}
		return bytes.Compare(nbx.nfts[j].Bytes(), nbx.nfts[i].Bytes()) < 0
	})
}

func (nbx NFTBox) Exists(id nft.NFTID) bool {
	if len(nbx.nfts) < 1 {
		return false
	}
	for _, n := range nbx.nfts {
		if id.Equal(n) {
			return true
		}
	}
	return false
}

func (nbx NFTBox) Get(id nft.NFTID) (*nft.NFTID, error) {
	for _, n := range nbx.nfts {
		if id.Equal(n) {
			return &n, nil
		}
	}
	return nil, errors.Errorf("nft not found in NFTBox, %q", id)
}

func (nbx *NFTBox) Append(n nft.NFTID) error {
	if err := n.IsValid(nil); err != nil {
		return err
	}
	if nbx.Exists(n) {
		return errors.Errorf("nft already exists in NFTBox, %q", n)
	}
	if uint64(len(nbx.nfts)) >= nft.MaxNFTIndex {
		return errors.Errorf("max nfts in collection, %q", n)
	}
	nbx.nfts = append(nbx.nfts, n)
	return nil
}

func (nbx *NFTBox) Remove(n nft.NFTID) error {
	if err := n.IsValid(nil); err != nil {
		return err
	}
	if !nbx.Exists(n) {
		return errors.Errorf("nft not found in NFTBox, %q", n)
	}
	for i := range nbx.nfts {
		if n.Equal(nbx.nfts[i]) {
			nbx.nfts[i] = nbx.nfts[len(nbx.nfts)-1]
			nbx.nfts = nbx.nfts[:len(nbx.nfts)-1]
			return nil
		}
	}
	return nil
}

func (nbx NFTBox) NFTs() []nft.NFTID {
	return nbx.nfts
}
