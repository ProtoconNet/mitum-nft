package state

import (
	"strings"

	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

type StateValueMerger struct {
	*common.BaseStateValueMerger
}

func NewStateValueMerger(height mitumbase.Height, key string, st mitumbase.State) *StateValueMerger {
	s := &StateValueMerger{
		BaseStateValueMerger: common.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewStateMergeValue(key string, stv mitumbase.StateValue) mitumbase.StateMergeValue {
	StateValueMergerFunc := func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
		return NewStateValueMerger(height, key, st)
	}

	return mitumbase.NewBaseStateMergeValue(
		key,
		stv,
		StateValueMergerFunc,
	)
}

var CollectionStateValueHint = hint.MustNewHint("collection-state-value-v0.0.1")

type CollectionStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewCollectionStateValue(design types.Design) CollectionStateValue {
	return CollectionStateValue{
		BaseHinter: hint.NewBaseHinter(CollectionStateValueHint),
		Design:     design,
	}
}

func (cs CollectionStateValue) Hint() hint.Hint {
	return cs.BaseHinter.Hint()
}

func (cs CollectionStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionStateValue")

	if err := cs.BaseHinter.IsValid(CollectionStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := cs.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (cs CollectionStateValue) HashBytes() []byte {
	return cs.Design.Bytes()
}

func StateCollectionValue(st mitumbase.State) (*types.Design, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("collection design not found in State")
	}

	d, ok := v.(CollectionStateValue)
	if !ok {
		return nil, errors.Errorf("invalid collection value found, %T", v)
	}

	return &d.Design, nil
}

var LastNFTIndexStateValueHint = hint.MustNewHint("collection-last-nft-index-state-value-v0.0.1")

type LastNFTIndexStateValue struct {
	hint.BaseHinter
	id uint64
}

func NewLastNFTIndexStateValue( /*collection currencytypes.ContractID,*/ id uint64) LastNFTIndexStateValue {
	return LastNFTIndexStateValue{
		BaseHinter: hint.NewBaseHinter(LastNFTIndexStateValueHint),
		id:         id,
	}
}

func (is LastNFTIndexStateValue) Hint() hint.Hint {
	return is.BaseHinter.Hint()
}

func (is LastNFTIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid LastNFTIndexStateValue")

	if err := is.BaseHinter.IsValid(LastNFTIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (is LastNFTIndexStateValue) HashBytes() []byte {
	return util.Uint64ToBytes(is.id)
}

func StateLastNFTIndexValue(st mitumbase.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(LastNFTIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.id, nil
}

var (
	NFTStateValueHint = hint.MustNewHint("nft-state-value-v0.0.1")
)

type NFTStateValue struct {
	hint.BaseHinter
	NFT types.NFT
}

func NewNFTStateValue(n types.NFT) NFTStateValue {
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

func StateNFTValue(st mitumbase.State) (*types.NFT, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("nft not found in State")
	}

	ns, ok := v.(NFTStateValue)
	if !ok {
		return nil, errors.Errorf("invalid nft value found, %T", v)
	}

	return &ns.NFT, nil
}

var NFTBoxStateValueHint = hint.MustNewHint("nft-box-state-value-v0.0.1")

type NFTBoxStateValue struct {
	hint.BaseHinter
	Box types.NFTBox
}

func NewNFTBoxStateValue(box types.NFTBox) NFTBoxStateValue {
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

func StateNFTBoxValue(st mitumbase.State) (types.NFTBox, error) {
	v := st.Value()
	if v == nil {
		return types.NFTBox{}, util.ErrNotFound.Errorf("nft box not found in State")
	}

	nb, ok := v.(NFTBoxStateValue)
	if !ok {
		return types.NFTBox{}, errors.Errorf("invalid nft box value found, %T", v)
	}

	return nb.Box, nil
}

var OperatorsBookStateValueHint = hint.MustNewHint("operators-book-state-value-v0.0.1")

type OperatorsBookStateValue struct {
	hint.BaseHinter
	Operators types.OperatorsBook
}

func NewOperatorsBookStateValue(operators types.OperatorsBook) OperatorsBookStateValue {
	return OperatorsBookStateValue{
		BaseHinter: hint.NewBaseHinter(OperatorsBookStateValueHint),
		Operators:  operators,
	}
}

func (ob OperatorsBookStateValue) Hint() hint.Hint {
	return ob.BaseHinter.Hint()
}

func (ob OperatorsBookStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid OperatorsBookStateValue")

	if err := ob.BaseHinter.IsValid(OperatorsBookStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ob.Operators.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ob OperatorsBookStateValue) HashBytes() []byte {
	return ob.Operators.Bytes()
}

func StateOperatorsBookValue(st mitumbase.State) (*types.OperatorsBook, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("operators book not found in State")
	}

	ob, ok := v.(OperatorsBookStateValue)
	if !ok {
		return nil, errors.Errorf("invalid operators book value found, %T", v)
	}

	return &ob.Operators, nil
}

// ParsedStateKey is the function that parses the state key.
// The length of state key is 4 or 5.
// In case of length 4 it forms as NFTPrefix:{contract}:{collection}:{Suffix}.
// In case of length 5 it forms as NFTPrefix:{contract}:{collection}:{key_value}:{Suffix}
func ParseStateKey(key string) ([]string, error) {
	parsedKey := strings.Split(key, ":")
	if parsedKey[0] != NFTPrefix[:len(NFTPrefix)-1] {
		return nil, errors.Errorf("State Key not include NFTPrefix, %s", parsedKey)
	}
	if len(parsedKey) < 3 {
		return nil, errors.Errorf("parsing State Key string failed, %s", parsedKey)
	} else {
		return parsedKey, nil
	}
}

func checkExistsState(
	key string,
	getState mitumbase.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return mitumbase.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkExistsStates(
	keys []string,
	getState mitumbase.GetStateFunc,
) error {
	for i := range keys {
		switch _, found, err := getState(keys[i]); {
		case err != nil:
			return err
		case !found:
			return mitumbase.NewBaseOperationProcessReasonError("state, %q does not exist", keys[i])
		}
	}
	return nil
}

func checkNotExistsState(
	key string,
	getState mitumbase.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return mitumbase.NewBaseOperationProcessReasonError("state, %q already exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState mitumbase.GetStateFunc,
) (mitumbase.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, mitumbase.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func existsStates(
	getState mitumbase.GetStateFunc,
	keys ...string,
) ([]mitumbase.State, error) {
	var states []mitumbase.State
	for i := range keys {
		switch st, found, err := getState(keys[i]); {
		case err != nil:
			return nil, err
		case !found:
			return nil, mitumbase.NewBaseOperationProcessReasonError("value of key does not exist, %s", keys[i])
		default:
			states = append(states, st)
		}
	}
	if len(keys) != len(states) {
		return nil, mitumbase.NewBaseOperationProcessReasonError("get multiple states failed")
	}
	return states, nil
}

func notExistsState(
	k,
	name string,
	getState mitumbase.GetStateFunc,
) (mitumbase.State, error) {
	var st mitumbase.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, mitumbase.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = mitumbase.NewBaseState(mitumbase.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currencytypes.CurrencyID, getStateFunc mitumbase.GetStateFunc) (currencytypes.CurrencyPolicy, error) {
	var policy currencytypes.CurrencyPolicy

	switch st, found, err := getStateFunc(statecurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return currencytypes.CurrencyPolicy{}, err
	case !found:
		return currencytypes.CurrencyPolicy{}, errors.Errorf("currency not found, %v", cid)
	default:
		design, ok := st.Value().(statecurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return currencytypes.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", st.Value())
		}
		policy = design.CurrencyDesign.Policy()
	}

	return policy, nil
}

func existsCollectionPolicy(contract mitumbase.Address, id currencytypes.ContractID, getStateFunc mitumbase.GetStateFunc) (types.CollectionPolicy, error) {
	var policy types.CollectionPolicy

	switch st, found, err := getStateFunc(NFTStateKey(contract, id, CollectionKey)); {
	case err != nil:
		return types.CollectionPolicy{}, err
	case !found:
		return types.CollectionPolicy{}, errors.Errorf("collection not found, %v", id)
	default:
		design, ok := st.Value().(CollectionStateValue)
		if !ok {
			return types.CollectionPolicy{}, errors.Errorf("expected CollectionDesignStateValue, not %T", st.Value())
		}
		p := design.Design.Policy()
		policy, ok = p.(types.CollectionPolicy)
		if !ok {
			return types.CollectionPolicy{}, errors.Errorf("expected CollectionPolicy, not %T", p)
		}
	}

	return policy, nil
}
