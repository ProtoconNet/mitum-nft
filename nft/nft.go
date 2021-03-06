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

var (
	MaxCreators     = 10
	MaxCopyrighters = 10
)

type NFT struct {
	hint.BaseHinter
	id           NFTID
	active       bool
	owner        base.Address
	hash         NFTHash
	uri          URI
	approved     base.Address
	creators     Signers
	copyrighters Signers
}

func NewNFT(id NFTID, active bool, owner base.Address, hash NFTHash, uri URI, approved base.Address, creators Signers, copyrighters Signers) NFT {
	return NFT{
		BaseHinter:   hint.NewBaseHinter(NFTHint),
		id:           id,
		active:       active,
		owner:        owner,
		hash:         hash,
		uri:          uri,
		approved:     approved,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func MustNewNFT(id NFTID, active bool, owner base.Address, hash NFTHash, uri URI, approved base.Address, creators Signers, copyrighters Signers) NFT {
	n := NewNFT(id, active, owner, hash, uri, approved, creators, copyrighters)

	if err := n.IsValid(nil); err != nil {
		panic(err)
	}

	return n
}

func (n NFT) Bytes() []byte {
	ba := make([]byte, 1)

	if n.active {
		ba[0] = 1
	} else {
		ba[0] = 0
	}

	return util.ConcatBytesSlice(
		n.id.Bytes(),
		ba,
		n.owner.Bytes(),
		n.hash.Bytes(),
		[]byte(n.uri.String()),
		n.approved.Bytes(),
		n.creators.Bytes(),
		n.copyrighters.Bytes(),
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
	if err := isvalid.Check(
		nil, false,
		n.id,
		n.owner,
		n.hash,
		n.uri,
		n.approved,
		n.creators,
		n.copyrighters,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid nft; %w", err)
	}

	if len(n.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	return nil
}

func (n NFT) ID() NFTID {
	return n.id
}

func (n NFT) Active() bool {
	return n.active
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

func (n NFT) Creators() Signers {
	return n.creators
}

func (n NFT) Copyrighters() Signers {
	return n.copyrighters
}

func (n NFT) Equal(cn NFT) bool {
	return n.ID().Equal(cn.ID())
}

func (n NFT) ExistsApproved() bool {
	return !n.approved.Equal(n.owner)
}
