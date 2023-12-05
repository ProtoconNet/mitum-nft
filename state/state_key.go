package state

import (
	"fmt"
	"strconv"
	"strings"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type StateKey int

const (
	NilKey = iota
	CollectionKey
	OperatorsKey
	LastIDXKey
	NFTBoxKey
	NFTKey
)

var (
	NFTPrefix                = "nft"
	StateKeyCollectionSuffix = "collection"
	StateKeyOperatorsSuffix  = "operators"
	StateKeyLastNFTIDXSuffix = "lastnftidx"
	StateKeyNFTBoxSuffix     = "nftbox"
	StateKeyNFTSuffix        = "nft"
)

func StateKeyNFTPrefix(addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s", NFTPrefix, addr.String())
}

func NFTStateKey(
	contract mitumbase.Address,
	keyType StateKey,
) string {
	prefix := StateKeyNFTPrefix(contract)
	var stateKey string
	switch keyType {
	case CollectionKey:
		stateKey = fmt.Sprintf("%s:%s", prefix, StateKeyCollectionSuffix)
	case LastIDXKey:
		stateKey = fmt.Sprintf("%s:%s", prefix, StateKeyLastNFTIDXSuffix)
	case NFTBoxKey:
		stateKey = fmt.Sprintf("%s:%s", prefix, StateKeyNFTBoxSuffix)
	}

	return stateKey
}

func StateKeyOperators(contract mitumbase.Address, addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyNFTPrefix(contract), addr.String(), StateKeyOperatorsSuffix)
}

func StateKeyNFT(contract mitumbase.Address, id uint64) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyNFTPrefix(contract), strconv.FormatUint(id, 10), StateKeyNFTSuffix)
}

func ParseNFTStateKey(key string) (StateKey, error) {
	if !strings.HasPrefix(key, NFTPrefix) {
		return NilKey, errors.Errorf("invalid NFT State Key")
	}
	switch {
	case strings.HasSuffix(key, StateKeyCollectionSuffix):
		return CollectionKey, nil
	case strings.HasSuffix(key, StateKeyNFTBoxSuffix):
		return NFTBoxKey, nil
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
