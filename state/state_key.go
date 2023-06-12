package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ProtoconNet/mitum-currency/v3/types"

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
	NFTPrefix                = "nft:"
	StateKeyCollectionSuffix = ":collection"
	StateKeyOperatorsSuffix  = ":operators"
	StateKeyLastNFTIDXSuffix = ":lastnftidx"
	StateKeyNFTBoxSuffix     = ":nftbox"
	StateKeyNFTSuffix        = ":nft"
)

func StateKeyNFTPrefix(addr mitumbase.Address, collectionID types.ContractID) string {
	return fmt.Sprintf("%s%s:%s", NFTPrefix, addr.String(), collectionID.String())
}

func NFTStateKey(
	contract mitumbase.Address,
	collectionID types.ContractID,
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
	case NFTBoxKey:
		stateKey = fmt.Sprintf("%s%s", prefix, StateKeyNFTBoxSuffix)
	}

	return stateKey
}

func StateKeyOperators(contract mitumbase.Address, collectionID types.ContractID, addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), addr.String(), StateKeyOperatorsSuffix)
}

func StateKeyNFT(contract mitumbase.Address, collectionID types.ContractID, id uint64) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), strconv.FormatUint(id, 10), StateKeyNFTSuffix)
}

func ParseNFTStateKey(key string) (StateKey, error) {
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
