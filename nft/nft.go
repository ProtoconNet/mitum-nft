package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

type NFTHash string

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

func (hs NFTHash) IsValid([]byte) error {
	return nil
}

var (
	NFTType   = hint.Type("mitum-nft-nft")
	NFTHint   = hint.NewHint(NFTType, "v0.0.1")
	NFTHinter = NFT{BaseHinter: hint.NewBaseHinter(NFTHint)}
)

type NFT struct {
	hint.BaseHinter
	id           NFTID
	owner        base.Address
	hash         NFTHash
	uri          URI
	approved     base.Address
	creators     []Signer
	copyrighters []Signer
}

func NewNFT(id NFTID, owner base.Address, hash NFTHash, uri URI, approved base.Address, creators []Signer, copyrighters []Signer) NFT {
	return NFT{
		BaseHinter:   hint.NewBaseHinter(NFTHint),
		id:           id,
		owner:        owner,
		hash:         hash,
		uri:          uri,
		approved:     approved,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func MustNewNFT(id NFTID, owner base.Address, hash NFTHash, uri URI, approved base.Address, creators []Signer, copyrighters []Signer) NFT {
	n := NewNFT(id, owner, hash, uri, approved, creators, copyrighters)

	if err := n.IsValid(nil); err != nil {
		panic(err)
	}

	return n
}

func (n NFT) Bytes() []byte {
	bcrs := [][]byte{}
	bcps := [][]byte{}

	for i := range n.creators {
		bcrs = append(bcrs, n.creators[i].Bytes())
	}

	for i := range n.copyrighters {
		bcps = append(bcrs, n.copyrighters[i].Bytes())
	}

	return util.ConcatBytesSlice(
		n.id.Bytes(),
		n.owner.Bytes(),
		n.hash.Bytes(),
		[]byte(n.uri.String()),
		n.approved.Bytes(),
		util.ConcatBytesSlice(bcrs...),
		util.ConcatBytesSlice(bcps...),
	)
}

func (NFT) Hint() hint.Hint {
	return NFTHint
}

func (n NFT) Hash() valuehash.Hash {
	return n.GenerateHash()
}

func (n NFT) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(n.Bytes())
}

func (n NFT) IsValid([]byte) error {
	if len(n.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	for i := range n.creators {
		if err := n.creators[i].IsValid(nil); err != nil {
			return err
		}
	}

	for i := range n.copyrighters {
		if err := n.copyrighters[i].IsValid(nil); err != nil {
			return err
		}
	}

	if len(n.approved.String()) > 0 {
		if err := n.approved.IsValid(nil); err != nil {
			return err
		}
	}

	if err := isvalid.Check(
		nil, false,
		n.id,
		n.hash,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid nft; %w", err)
	}

	return nil
}

func (n NFT) ID() NFTID {
	return n.id
}

func (n NFT) Owner() base.Address {
	return n.owner
}

func (n NFT) NftHash() NFTHash {
	return n.hash
}

func (n NFT) Uri() URI {
	return n.uri
}

func (n NFT) Approved() base.Address {
	return n.approved
}

func (n NFT) Creators() []Signer {
	return n.creators
}

func (n NFT) Copyrighters() []Signer {
	return n.copyrighters
}

func (n NFT) Equal(cn NFT) bool {
	return n.ID().Equal(cn.ID())
}
