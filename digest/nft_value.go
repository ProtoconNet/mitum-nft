package digest

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	NFTValueType = hint.Type("mitum-nft-value")
	NFTValueHint = hint.NewHint(NFTValueType, "v0.0.1")
)

type NFTValue struct {
	nft    nft.NFT
	height base.Height
}

func NewNFTValue(
	doc nft.NFT,
	height base.Height,
) NFTValue {
	return NFTValue{
		nft:    doc,
		height: height,
	}
}

func (NFTValue) Hint() hint.Hint {
	return NFTValueHint
}

func (n NFTValue) NFT() nft.NFT {
	return n.nft
}

func (n NFTValue) Height() base.Height {
	return n.height
}

func DecodeNFT(b []byte, enc encoder.Encoder) (nft.NFT, error) {
	if i, err := enc.Decode(b); err != nil {
		return nft.NFT{}, err
	} else if i == nil {
		return nft.NFT{}, nil
	} else if v, ok := i.(nft.NFT); !ok {
		return nft.NFT{}, util.WrongTypeError.Errorf("not NFT; type=%T", i)
	} else {
		return v, nil
	}
}
