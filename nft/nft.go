package nft

import (
	"fmt"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var BLACKHOLE_ZERO = currency.NewAddress("blackhole-0")

func isCollectionEqual(c1 extensioncurrency.ContractID, c2 extensioncurrency.ContractID) bool {
	return c1.String() == c2.String()
}

var (
	NFTIDType   = hint.Type("mitum-nft-nft-id")
	NFTIDHint   = hint.NewHint(NFTIDType, "v0.0.1")
	NFTIDHinter = NFTID{BaseHinter: hint.NewBaseHinter(NFTIDHint)}
)

var MaxNFTsInCollection = 10000

type NFTID struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	idx        uint
}

func NewNFTID(collection extensioncurrency.ContractID, idx uint) NFTID {
	return NFTID{
		BaseHinter: hint.NewBaseHinter(NFTIDHint),
		collection: collection,
		idx:        idx,
	}
}

func MustNewNFTID(collection extensioncurrency.ContractID, idx uint) NFTID {
	id := NewNFTID(collection, idx)

	if err := id.IsValid(nil); err != nil {
		panic(err)
	}

	return id
}

func (nid NFTID) Bytes() []byte {
	return util.ConcatBytesSlice(
		nid.collection.Bytes(),
		util.UintToBytes(nid.idx),
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

func (nid NFTID) IsValid([]byte) error {
	if nid.idx > uint(MaxNFTsInCollection) {
		return isvalid.InvalidError.Errorf("nft idx over max value; %d < %d", MaxNFTsInCollection, nid.idx)
	}

	if err := isvalid.Check(nil, false,
		nid.BaseHinter,
		nid.collection,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid nft id; %w", err)
	}

	return nil
}

func (nid NFTID) Symbol() extensioncurrency.ContractID {
	return nid.collection
}

func (nid NFTID) Idx() uint {
	return nid.idx
}

func (nid NFTID) String() string {
	return fmt.Sprintf("%s-%d)", nid.collection.String(), nid.idx)
}

func (nid NFTID) Equal(cnid NFTID) bool {
	return isCollectionEqual(nid.collection, cnid.collection) && nid.idx == cnid.idx
}

type NFTHash string

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

func (hs NFTHash) IsValid([]byte) error {
	if len(hs) == 0 {
		return isvalid.InvalidError.Errorf("empty nft hash")
	}

	return nil
}

var (
	NFTType   = hint.Type("mitum-nft-nft")
	NFTHint   = hint.NewHint(NFTType, "v0.0.1")
	NFTHinter = NFT{BaseHinter: hint.NewBaseHinter(NFTHint)}
)

type NFT struct {
	hint.BaseHinter
	id          NFTID
	owner       base.Address
	hash        NFTHash
	uri         URI
	approved    base.Address
	copyrighter base.Address
}

func NewNFT(id NFTID, owner base.Address, hash NFTHash, uri URI, approved base.Address, copyrighter base.Address) NFT {
	return NFT{
		BaseHinter:  hint.NewBaseHinter(NFTHint),
		id:          id,
		owner:       owner,
		hash:        hash,
		uri:         uri,
		approved:    approved,
		copyrighter: copyrighter,
	}
}

func MustNewNFT(id NFTID, owner base.Address, hash NFTHash, uri URI, approved base.Address, copyrighter base.Address) NFT {
	nft := NewNFT(id, owner, hash, uri, approved, copyrighter)

	if err := nft.IsValid(nil); err != nil {
		panic(err)
	}

	return nft
}

func (nft NFT) Bytes() []byte {
	return util.ConcatBytesSlice(
		nft.id.Bytes(),
		nft.owner.Bytes(),
		nft.hash.Bytes(),
		[]byte(nft.uri.String()),
		nft.approved.Bytes(),
		nft.copyrighter.Bytes(),
	)
}

func (nft NFT) Hint() hint.Hint {
	return NFTHint
}

func (nft NFT) Hash() valuehash.Hash {
	return nft.GenerateHash()
}

func (nft NFT) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(nft.Bytes())
}

func (nft NFT) IsValid([]byte) error {
	if len(nft.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	if len(nft.copyrighter.String()) > 1 {
		if err := nft.copyrighter.IsValid(nil); err != nil {
			return err
		}
	}

	if err := isvalid.Check(
		nil, false,
		nft.id,
		nft.owner,
		nft.hash,
		nft.approved,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid nft; %w", err)
	}
	return nil
}

func (nft NFT) ID() NFTID {
	return nft.id
}

func (nft NFT) Owner() base.Address {
	return nft.owner
}

func (nft NFT) NftHash() NFTHash {
	return nft.hash
}

func (nft NFT) Uri() URI {
	return nft.uri
}

func (nft NFT) Approved() base.Address {
	return nft.approved
}

func (nft NFT) Copyrighter() base.Address {
	return nft.copyrighter
}

func (nft NFT) Equal(cnft NFT) bool {
	return nft.ID().Equal(cnft.ID())
}
