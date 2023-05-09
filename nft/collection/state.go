package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	CollectionDesignStateValueHint = hint.MustNewHint("collection-design-state-value-v0.0.1")
	StateKeyCollectionDesignPrefix = "collection:"
)

type CollectionDesignStateValue struct {
	hint.BaseHinter
	CollectionDesign CollectionDesign
}

func NewCollectionDesignStateValue(design CollectionDesign) CollectionDesignStateValue {
	return CollectionDesignStateValue{
		BaseHinter:       hint.NewBaseHinter(CollectionDesignStateValueHint),
		CollectionDesign: design,
	}
}

func (cs CollectionDesignStateValue) Hint() hint.Hint {
	return cs.BaseHinter.Hint()
}

func (cs CollectionDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionDesignStateValue")

	if err := cs.BaseHinter.IsValid(CollectionDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := cs.CollectionDesign.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (cs CollectionDesignStateValue) HashBytes() []byte {
	return cs.CollectionDesign.Bytes()
}

func StateCollectionDesignValue(st base.State) (CollectionDesign, error) {
	v := st.Value()
	if v == nil {
		return CollectionDesign{}, util.ErrNotFound.Errorf("collection design not found in State")
	}

	d, ok := v.(CollectionDesignStateValue)
	if !ok {
		return CollectionDesign{}, errors.Errorf("invalid collection design value found, %T", v)
	}

	return d.CollectionDesign, nil
}

type CollectionDesignStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewCollectionDesignStateValueMerger(height base.Height, key string, st base.State) *CollectionDesignStateValueMerger {
	s := &CollectionDesignStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewCollectionDesignStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewCollectionDesignStateValueMerger(height, key, st)
		},
	)
}

var (
	CollectionLastNFTIndexStateValueHint = hint.MustNewHint("collection-last-nft-index-state-value-v0.0.1")
	StateKeyCollectionLastNFTIndexSuffix = ":collectionidx"
)

type CollectionLastNFTIndexStateValue struct {
	hint.BaseHinter
	// Collection extensioncurrency.ContractID
	Index uint64
}

func NewCollectionLastNFTIndexStateValue( /*collection extensioncurrency.ContractID,*/ index uint64) CollectionLastNFTIndexStateValue {
	return CollectionLastNFTIndexStateValue{
		BaseHinter: hint.NewBaseHinter(CollectionLastNFTIndexStateValueHint),
		// Collection: collection,
		Index: index,
	}
}

func (is CollectionLastNFTIndexStateValue) Hint() hint.Hint {
	return is.BaseHinter.Hint()
}

func (is CollectionLastNFTIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionLastNFTIndexStateValue")

	if err := is.BaseHinter.IsValid(CollectionLastNFTIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	// if err := is.Collection.IsValid(nil); err != nil {
	// 	return e.Wrap(err)
	// }

	return nil
}

func (is CollectionLastNFTIndexStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice( /*is.Collection.Bytes(), */ util.Uint64ToBytes(is.Index))
}

func StateCollectionLastNFTIndexValue(st base.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(CollectionLastNFTIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.Index, nil
}

type CollectionLastNFTIndexStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewCollectionLastNFTIndexStateValueMerger(height base.Height, key string, st base.State) *CollectionLastNFTIndexStateValueMerger {
	s := &CollectionLastNFTIndexStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewCollectionLastNFTIndexStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewCollectionLastNFTIndexStateValueMerger(height, key, st)
		},
	)
}

var (
	NFTStateValueHint = hint.MustNewHint("nft-state-value-v0.0.1")
)

type NFTStateValue struct {
	hint.BaseHinter
	NFT nft.NFT
}

func NewNFTStateValue(n nft.NFT) NFTStateValue {
	return NFTStateValue{
		BaseHinter: hint.NewBaseHinter(NFTStateValueHint),
		NFT:        n,
	}
}

func (ns NFTStateValue) Hint() hint.Hint {
	return ns.BaseHinter.Hint()
}

func (ns NFTStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NFTStateValue")

	if err := ns.BaseHinter.IsValid(NFTStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ns.NFT.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ns NFTStateValue) HashBytes() []byte {
	return ns.NFT.Bytes()
}

func StateNFTValue(st base.State) (nft.NFT, error) {
	v := st.Value()
	if v == nil {
		return nft.NFT{}, util.ErrNotFound.Errorf("nft not found in State")
	}

	ns, ok := v.(NFTStateValue)
	if !ok {
		return nft.NFT{}, errors.Errorf("invalid nft value found, %T", v)
	}

	return ns.NFT, nil
}

type NFTStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewNFTStateValueMerger(height base.Height, key string, st base.State) *NFTStateValueMerger {
	s := &NFTStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewNFTStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewNFTStateValueMerger(height, key, st)
		},
	)
}

var (
	NFTBoxStateValueHint = hint.MustNewHint("nft-box-state-value-v0.0.1")
	StateKeyNFTBoxSuffix = ":nftbox"
)

type NFTBoxStateValue struct {
	hint.BaseHinter
	Box NFTBox
}

func NewNFTBoxStateValue(box NFTBox) NFTBoxStateValue {
	return NFTBoxStateValue{
		BaseHinter: hint.NewBaseHinter(NFTBoxStateValueHint),
		Box:        box,
	}
}

func (nb NFTBoxStateValue) Hint() hint.Hint {
	return nb.BaseHinter.Hint()
}

func (nb NFTBoxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NFTBoxStateValue")

	if err := nb.BaseHinter.IsValid(NFTBoxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := nb.Box.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (nb NFTBoxStateValue) HashBytes() []byte {
	return nb.Box.Bytes()
}

func StateNFTBoxValue(st base.State) (NFTBox, error) {
	v := st.Value()
	if v == nil {
		return NFTBox{}, util.ErrNotFound.Errorf("nft box not found in State")
	}

	nb, ok := v.(NFTBoxStateValue)
	if !ok {
		return NFTBox{}, errors.Errorf("invalid nft box value found, %T", v)
	}

	return nb.Box, nil
}

type NFTBoxStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewNFTBoxStateValueMerger(height base.Height, key string, st base.State) *NFTBoxStateValueMerger {
	s := &NFTBoxStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewNFTBoxStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewNFTBoxStateValueMerger(height, key, st)
		},
	)
}

var (
	AgentBoxStateValueHint = hint.MustNewHint("agent-box-state-value-v0.0.1")
	StateKeyAgentBoxSuffix = ":agentbox"
)

type AgentBoxStateValue struct {
	hint.BaseHinter
	Box AgentBox
}

func NewAgentBoxStateValue(box AgentBox) AgentBoxStateValue {
	return AgentBoxStateValue{
		BaseHinter: hint.NewBaseHinter(AgentBoxStateValueHint),
		Box:        box,
	}
}

func (ab AgentBoxStateValue) Hint() hint.Hint {
	return ab.BaseHinter.Hint()
}

func (ab AgentBoxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid AgentBoxStateValue")

	if err := ab.BaseHinter.IsValid(AgentBoxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ab.Box.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ab AgentBoxStateValue) HashBytes() []byte {
	return ab.Box.Bytes()
}

func StateAgentBoxValue(st base.State) (AgentBox, error) {
	v := st.Value()
	if v == nil {
		return AgentBox{}, util.ErrNotFound.Errorf("agent box not found in State")
	}

	ab, ok := v.(AgentBoxStateValue)
	if !ok {
		return AgentBox{}, errors.Errorf("invalid agent box value found, %T", v)
	}

	return ab.Box, nil
}

type AgentBoxStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewAgentBoxStateValueMerger(height base.Height, key string, st base.State) *AgentBoxStateValueMerger {
	s := &AgentBoxStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewAgentBoxStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewNFTBoxStateValueMerger(height, key, st)
		},
	)
}

func checkExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return base.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkExistsStates(
	keys []string,
	getState base.GetStateFunc,
) error {
	for i := range keys {
		switch _, found, err := getState(keys[i]); {
		case err != nil:
			return err
		case !found:
			return base.NewBaseOperationProcessReasonError("state, %q does not exist", keys[i])
		}
	}
	return nil
}

func checkNotExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return base.NewBaseOperationProcessReasonError("state, %q already exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func existsStates(
	getState base.GetStateFunc,
	keys ...string,
) ([]base.State, error) {
	var states []base.State
	for i := range keys {
		switch st, found, err := getState(keys[i]); {
		case err != nil:
			return nil, err
		case !found:
			return nil, base.NewBaseOperationProcessReasonError("value of key does not exist, %s", keys[i])
		default:
			states = append(states, st)
		}
	}
	if len(keys) != len(states) {
		return nil, base.NewBaseOperationProcessReasonError("get multiple states failed")
	}
	return states, nil
}

func notExistsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	var st base.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, base.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = base.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currency.CurrencyID, getStateFunc base.GetStateFunc) (extensioncurrency.CurrencyPolicy, error) {
	var policy extensioncurrency.CurrencyPolicy

	switch st, found, err := getStateFunc(extensioncurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return extensioncurrency.CurrencyPolicy{}, err
	case !found:
		return extensioncurrency.CurrencyPolicy{}, errors.Errorf("currency not found, %v", cid)
	default:
		design, ok := st.Value().(extensioncurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return extensioncurrency.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", st.Value())
		}
		policy = design.CurrencyDesign.Policy()
	}

	return policy, nil
}

func existsCollectionPolicy(contract base.Address, id extensioncurrency.ContractID, getStateFunc base.GetStateFunc) (CollectionPolicy, error) {
	var policy CollectionPolicy

	switch st, found, err := getStateFunc(NFTStateKey(contract, id, CollectionKey)); {
	case err != nil:
		return CollectionPolicy{}, err
	case !found:
		return CollectionPolicy{}, errors.Errorf("collection not found, %v", id)
	default:
		design, ok := st.Value().(CollectionDesignStateValue)
		if !ok {
			return CollectionPolicy{}, errors.Errorf("expected CollectionDesignStateValue, not %T", st.Value())
		}
		p := design.CollectionDesign.Policy()
		policy, ok = p.(CollectionPolicy)
		if !ok {
			return CollectionPolicy{}, errors.Errorf("expected CollectionPolicy, not %T", p)
		}
	}

	return policy, nil
}
