package collection

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type StateKey int

const (
	NilKey = iota
	CollectionKey
	OperatorsKey
	LastIDXKey
	NFTsKey
	NFTKey
)

var (
	NFTPrefix                = "nft:"
	StateKeyCollectionSuffix = ":collection"
	StateKeyOperatorsSuffix  = ":operators"
	StateKeyLastNFTIDXSuffix = ":lastnftidx"
	StateKeyNFTsSuffix       = ":nfts"
	StateKeyNFTSuffix        = ":nft"
)

func StateKeyNFTPrefix(addr base.Address, collectionID extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s:%s", NFTPrefix, addr.String(), collectionID.String())
}

func NFTStateKey(
	contract base.Address,
	collectionID extensioncurrency.ContractID,
	keyType StateKey,
) string {
	prefix := StateKeyNFTPrefix(contract, collectionID)
	var stateKey string
	switch keyType {
	case CollectionKey:
		stateKey = fmt.Sprintf("%s%s", prefix, StateKeyCollectionSuffix)
	case OperatorsKey:
		stateKey = fmt.Sprintf("%s%s", prefix, StateKeyOperatorsSuffix)
	case LastIDXKey:
		stateKey = fmt.Sprintf("%s%s", prefix, StateKeyLastNFTIDXSuffix)
	case NFTsKey:
		stateKey = fmt.Sprintf("%s%s", prefix, StateKeyNFTsSuffix)
	}

	return stateKey
}

func StateKeyOperators(contract base.Address, collectionID extensioncurrency.ContractID, addr base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), addr.String(), StateKeyOperatorsSuffix)
}

func StateKeyNFT(contract base.Address, collectionID extensioncurrency.ContractID, id nft.NFTID) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), id.String(), StateKeyNFTSuffix)
}

func ParseNFTStateKey(key string) (StateKey, error) {
	switch {
	case strings.HasSuffix(key, StateKeyCollectionSuffix):
		return CollectionKey, nil
	case strings.HasSuffix(key, StateKeyNFTsSuffix):
		return NFTsKey, nil
	case strings.HasSuffix(key, StateKeyNFTSuffix):
		return NFTKey, nil
	case strings.HasSuffix(key, StateKeyLastNFTIDXSuffix):
		return LastIDXKey, nil
	case strings.HasSuffix(key, StateKeyOperatorsSuffix):
		return OperatorsKey, nil
	default:
		return NilKey, errors.Errorf("invalid NFT State Key")
	}
}
