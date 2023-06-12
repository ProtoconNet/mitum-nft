package types

import (
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MaxNFTHashLength = 1024
var MaxNFTIndex uint64 = 10000

type NFTHash string

func (hs NFTHash) IsValid([]byte) error {
	if l := len(hs); l > MaxNFTHashLength {
		return util.ErrInvalid.Errorf("nft hash length over max, %d > %d", l, MaxNFTHashLength)
	}

	if hs != "" && strings.TrimSpace(string(hs)) == "" {
		return util.ErrInvalid.Errorf("empty nft hash")
	}

	return nil
}

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

var NFTHint = hint.MustNewHint("mitum-nft-nft-v0.0.1")

var MaxCreators = 10

type NFT struct {
	hint.BaseHinter
	id       uint64
	active   bool
	owner    base.Address
	hash     NFTHash
	uri      URI
	approved base.Address
	creators Signers
}

func NewNFT(
	id uint64,
	active bool,
	owner base.Address,
	hash NFTHash,
	uri URI,
	approved base.Address,
	creators Signers,
) NFT {
	return NFT{
		BaseHinter: hint.NewBaseHinter(NFTHint),
		id:         id,
		active:     active,
		owner:      owner,
		hash:       hash,
		uri:        uri,
		approved:   approved,
		creators:   creators,
	}
}

func (n NFT) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		n.owner,
		n.hash,
		n.uri,
		n.approved,
		n.creators,
	); err != nil {
		return err
	}

	if n.uri == "" {
		return util.ErrInvalid.Errorf("empty uri")
	}

	return nil
}

func (n NFT) Bytes() []byte {
	ba := make([]byte, 1)

	if n.active {
		ba[0] = 1
	} else {
		ba[0] = 0
	}

	return util.ConcatBytesSlice(
		util.Uint64ToBytes(n.id),
		ba,
		n.owner.Bytes(),
		n.hash.Bytes(),
		[]byte(n.uri.String()),
		n.approved.Bytes(),
		n.creators.Bytes(),
	)
}

func (n NFT) ID() uint64 {
	return n.id
}

func (n NFT) Active() bool {
	return n.active
}

func (n NFT) Owner() base.Address {
	return n.owner
}

func (n NFT) NFTHash() NFTHash {
	return n.hash
}

func (n NFT) URI() URI {
	return n.uri
}

func (n NFT) Approved() base.Address {
	return n.approved
}

func (n NFT) Creators() Signers {
	return n.creators
}

func (n NFT) Addresses() []base.Address {
	var as []base.Address
	copy(as, n.Creators().Addresses())
	for i, a := range as {
		if n.approved != a {
			break
		}
		if i == (len(as) - 1) {
			as = append(as, n.approved)
		}
	}
	as = append(as)

	for i, a := range as {
		if n.owner != a {
			break
		}
		if i == (len(as) - 1) {
			as = append(as, n.owner)
		}
	}
	as = append(as)

	return as
}

func (n NFT) Equal(cn NFT) bool {
	if !(n.ID() == cn.ID()) {
		return false
	}

	if n.Active() != cn.Active() {
		return false
	}

	if !n.Owner().Equal(cn.Owner()) {
		return false
	}

	if n.NFTHash() != cn.NFTHash() {
		return false
	}

	if n.URI() != cn.URI() {
		return false
	}

	if !n.Approved().Equal(cn.Approved()) {
		return false
	}

	if !n.Creators().Equal(cn.Creators()) {
		return false
	}

	return n.ID() == cn.ID()
}

func (n NFT) ExistsApproved() bool {
	return !n.approved.Equal(n.owner)
}
