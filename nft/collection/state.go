package collection

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var (
	StateKeyCollectionPrefix = "collection:"
)

var (
	StateKeyAgentsSuffix            = ":agents"
	StateKeyCollectionLastIDXSuffix = ":collectionidx"
	StateKeyNFTsSuffix              = ":nfts"
	StateKeyNFTSuffix               = ":nft"
)

func StateKeyAgents(addr base.Address) string {
	return fmt.Sprintf("%s%s", addr.String(), StateKeyAgentsSuffix)
}

func IsStateAgentKey(key string) bool {
	return strings.HasSuffix(key, StateKeyAgentsSuffix)
}

func StateAgentsValue(st state.State) (AgentBox, error) {
	value := st.Value()
	if value == nil {
		return AgentBox{}, util.NotFoundError.Errorf("agent box not found in State")
	}

	if box, ok := value.Interface().(AgentBox); !ok {
		return AgentBox{}, errors.Errorf("invalid agent box value found; %T", value.Interface())
	} else {
		return box, nil
	}
}

func SetStateAgentsValue(st state.State, box AgentBox) (state.State, error) {
	if vbox, err := state.NewHintedValue(box); err != nil {
		return nil, err
	} else {
		return st.SetValue(vbox)
	}
}

func StateKeyCollection(id extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyCollectionPrefix, id.String())
}

func IsStateCollectionKey(key string) bool {
	return strings.HasPrefix(key, StateKeyCollectionPrefix)
}

func StateCollectionValue(st state.State) (nft.Design, error) {
	value := st.Value()
	if value == nil {
		return nft.Design{}, util.NotFoundError.Errorf("design not found in State")
	}

	if design, ok := value.Interface().(nft.Design); !ok {
		return nft.Design{}, errors.Errorf("invalid design value found; %T", value.Interface())
	} else {
		return design, nil
	}
}

func SetStateCollectionValue(st state.State, design nft.Design) (state.State, error) {
	if vdesign, err := state.NewHintedValue(design); err != nil {
		return nil, err
	} else {
		return st.SetValue(vdesign)
	}
}

func StateKeyNFTs(addr base.Address) string {
	return fmt.Sprintf("%s%s", addr.String(), StateKeyNFTsSuffix)
}

func IsStateNFTsKey(key string) bool {
	return strings.HasSuffix(key, StateKeyNFTsSuffix)
}

func StateNFTsValue(st state.State) (NFTBox, error) {
	value := st.Value()
	if value == nil {
		return NFTBox{}, util.NotFoundError.Errorf("nft box not found in State")
	}

	if box, ok := value.Interface().(NFTBox); !ok {
		return NFTBox{}, errors.Errorf("invalid nft box value found; %T", value.Interface())
	} else {
		return box, nil
	}
}

func SetStateNFTsValue(st state.State, box NFTBox) (state.State, error) {
	if vbox, err := state.NewHintedValue(box); err != nil {
		return nil, err
	} else {
		return st.SetValue(vbox)
	}
}

func StateKeyNFT(id nft.NFTID) string {
	return fmt.Sprintf("%s%s", id.String(), StateKeyNFTSuffix)
}

func IsStateNFTKey(key string) bool {
	return strings.HasSuffix(key, StateKeyNFTSuffix)
}

func StateNFTValue(st state.State) (nft.NFT, error) {
	value := st.Value()
	if value == nil {
		return nft.NFT{}, util.NotFoundError.Errorf("nft not found in State")
	}

	if n, ok := value.Interface().(nft.NFT); !ok {
		return nft.NFT{}, errors.Errorf("invalid nft value found; %T", value.Interface())
	} else {
		return n, nil
	}
}

func SetStateNFTValue(st state.State, n nft.NFT) (state.State, error) {
	if vn, err := state.NewHintedValue(n); err != nil {
		return nil, err
	} else {
		return st.SetValue(vn)
	}
}

func StateKeyCollectionLastIDX(id extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", id.String(), StateKeyCollectionLastIDXSuffix)
}

func IsStateCollectionLastIDXKey(key string) bool {
	return strings.HasSuffix(key, StateKeyCollectionLastIDXSuffix)
}

func StateCollectionLastIDXValue(st state.State) (currency.Big, error) {
	value := st.Value()
	if value == nil {
		return currency.Big{}, util.NotFoundError.Errorf("collection idx not found in State")
	}

	if idx, ok := value.Interface().(currency.Big); !ok {
		return currency.Big{}, errors.Errorf("invalid collection idx value found; %T", value.Interface())
	} else {
		return idx, nil
	}
}

func SetStateCollectionLastIDXValue(st state.State, idx currency.Big) (state.State, error) {
	if vidx, err := state.NewNumberValue(idx); err != nil {
		return nil, err
	} else {
		return st.SetValue(vidx)
	}
}

func checkExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return nil
	default:
		return operation.NewBaseReasonError("state, %q already exists", key)
	}
}

func existsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, operation.NewBaseReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, operation.NewBaseReasonError("%s already exists", name)
	default:
		return st, nil
	}
}
