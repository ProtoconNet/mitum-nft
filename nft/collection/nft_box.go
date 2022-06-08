package collection

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	NFTBoxType   = hint.Type("mitum-nft-nft-box")
	NFTBoxHint   = hint.NewHint(NFTBoxType, "v0.0.1")
	NFTBoxHinter = NFTBox{BaseHinter: hint.NewBaseHinter(NFTBoxHint)}
)

type NFTBox struct {
	hint.BaseHinter
	nfts []nft.NFTID
}

func NewNFTBox(nfts []nft.NFTID) NFTBox {
	if nfts == nil {
		return NFTBox{BaseHinter: hint.NewBaseHinter(NFTBoxHint), nfts: []nft.NFTID{}}
	}
	return NFTBox{BaseHinter: hint.NewBaseHinter(NFTBoxHint), nfts: nfts}
}

func (nbx NFTBox) Bytes() []byte {
	bs := make([][]byte, len(nbx.nfts))
	for i := range nbx.nfts {
		bs[i] = nbx.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (nbx NFTBox) Hint() hint.Hint {
	return NFTBoxHint
}

func (nbx NFTBox) Hash() valuehash.Hash {
	return nbx.GenerateHash()
}

func (nbx NFTBox) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(nbx.Bytes())
}

func (nbx NFTBox) IsEmpty() bool {
	return len(nbx.nfts) < 1
}

func (nbx NFTBox) IsValid([]byte) error {
	for i := range nbx.nfts {
		if err := nbx.nfts[i].IsValid(nil); err != nil {
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
	for i := range nbx.nfts {
		if id.Equal(nbx.nfts[i]) {
			return true
		}
	}
	return false
}

func (nbx NFTBox) Get(id nft.NFTID) (nft.NFTID, error) {
	for i := range nbx.nfts {
		if id.Equal(nbx.nfts[i]) {
			return nbx.nfts[i], nil
		}
	}
	return nft.NFTID{}, errors.Errorf("nft not found in owner's nft box; %v", id)
}

func (nbx *NFTBox) Append(n nft.NFTID) error {
	if err := n.IsValid(nil); err != nil {
		return err
	}
	if nbx.Exists(n) {
		return errors.Errorf("nft %v already exists in nft box", n)
	}
	if len(nbx.nfts) >= nft.MaxNFTsInCollection {
		return errors.Errorf("max nfts in collection; %v", n)
	}
	nbx.nfts = append(nbx.nfts, n)
	return nil
}

func (nbx *NFTBox) Remove(n nft.NFTID) error {
	if err := n.IsValid(nil); err != nil {
		return err
	}
	if !nbx.Exists(n) {
		return errors.Errorf("nft %v not found in nft box", n)
	}
	for i := range nbx.nfts {
		if n.Equal(nbx.nfts[i]) {
			nbx.nfts[i] = nbx.nfts[len(nbx.nfts)-1]
			nbx.nfts[len(nbx.nfts)-1] = nft.NFTID{}
			nbx.nfts = nbx.nfts[:len(nbx.nfts)-1]
			return nil
		}
	}
	return nil
}

func (nbx NFTBox) NFTs() []nft.NFTID {
	return nbx.nfts
}

type NFTBoxJSONPacker struct {
	jsonenc.HintedHead
	NS []nft.NFTID `json:"nfts"`
}

func (nbx NFTBox) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTBoxJSONPacker{
		HintedHead: jsonenc.NewHintedHead(nbx.Hint()),
		NS:         nbx.nfts,
	})
}

type NFTBoxJSONUnpacker struct {
	NS json.RawMessage `json:"nfts"`
}

func (nbx *NFTBox) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var un NFTBoxJSONUnpacker
	if err := enc.Unmarshal(b, &un); err != nil {
		return err
	}

	return nbx.unpack(enc, un.NS)
}

type NFTBoxBSONPacker struct {
	AG []nft.NFTID `bson:"nfts"`
}

func (nbx NFTBox) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(nbx.Hint()),
		bson.M{
			"nfts": nbx.nfts,
		}),
	)
}

type NFTBoxBSONUnpacker struct {
	NS bson.Raw `bson:"nfts"`
}

func (nbx *NFTBox) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var un NFTBoxBSONUnpacker
	if err := bsonenc.Unmarshal(b, &un); err != nil {
		return err
	}

	return nbx.unpack(enc, un.NS)
}

func (nbx *NFTBox) unpack(
	enc encoder.Encoder,
	bNFTs []byte,
) error {

	hNFTs, err := enc.DecodeSlice(bNFTs)
	if err != nil {
		return err
	}

	nfts := make([]nft.NFTID, len(hNFTs))
	for i := range hNFTs {
		j, ok := hNFTs[i].(nft.NFTID)
		if !ok {
			return util.WrongTypeError.Errorf("not NFTID; %T", hNFTs[i])
		}

		nfts[i] = j
	}

	nbx.nfts = nfts

	return nil
}
