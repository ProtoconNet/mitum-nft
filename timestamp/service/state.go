package service

import (
	"fmt"
	"strconv"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	StateKeyTimeStampPrefix     = "timestamp:"
	ServiceDesignStateValueHint = hint.MustNewHint("mitum-timestamp-service-design-state-value-v0.0.1")
	StateKeyServiceDesignSuffix = ":service"
)

func StateKeyTimeStampService(addr base.Address, sid extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s:%s", StateKeyTimeStampPrefix, addr.String(), sid)
}

type ServiceDesignStateValue struct {
	hint.BaseHinter
	Design timestamp.Design
}

func NewServiceDesignStateValue(design timestamp.Design) ServiceDesignStateValue {
	return ServiceDesignStateValue{
		BaseHinter: hint.NewBaseHinter(ServiceDesignStateValueHint),
		Design:     design,
	}
}

func (sd ServiceDesignStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd ServiceDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ServiceDesignStateValue")

	if err := sd.BaseHinter.IsValid(ServiceDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd ServiceDesignStateValue) HashBytes() []byte {
	return sd.Design.Bytes()
}

func StateServiceDesignValue(st base.State) (timestamp.Design, error) {
	v := st.Value()
	if v == nil {
		return timestamp.Design{}, util.ErrNotFound.Errorf("service design not found in State")
	}

	d, ok := v.(ServiceDesignStateValue)
	if !ok {
		return timestamp.Design{}, errors.Errorf("invalid service design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateServiceDesignKey(key string) bool {
	return strings.HasSuffix(key, StateKeyServiceDesignSuffix)
}

func StateKeyServiceDesign(addr base.Address, sid extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyTimeStampService(addr, sid), StateKeyServiceDesignSuffix)
}

type StateValueMerger struct {
	*currency.BaseStateValueMerger
}

func NewStateValueMerger(height base.Height, key string, st base.State) *StateValueMerger {
	s := &StateValueMerger{
		BaseStateValueMerger: currency.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	StateValueMergerFunc := func(height base.Height, st base.State) base.StateValueMerger {
		return NewStateValueMerger(height, key, st)
	}

	return currency.NewBaseStateMergeValue(
		key,
		stv,
		StateValueMergerFunc,
	)
}

var (
	TimeStampLastIndexStateValueHint = hint.MustNewHint("mitum-timestamp-last-index-state-value-v0.0.1")
	StateKeyProjectLastIndexSuffix   = ":timestampidx"
)

type TimeStampLastIndexStateValue struct {
	hint.BaseHinter
	ProjectID string
	Index     uint64
}

func NewTimeStampLastIndexStateValue(pid string, index uint64) TimeStampLastIndexStateValue {
	return TimeStampLastIndexStateValue{
		BaseHinter: hint.NewBaseHinter(TimeStampLastIndexStateValueHint),
		ProjectID:  pid,
		Index:      index,
	}
}

func (ti TimeStampLastIndexStateValue) Hint() hint.Hint {
	return ti.BaseHinter.Hint()
}

func (ti TimeStampLastIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TimeStampLastIndexStateValue")

	if err := ti.BaseHinter.IsValid(TimeStampLastIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(ti.ProjectID) < 1 || len(ti.ProjectID) > timestamp.MaxProjectIDLen {
		return errors.Errorf("invalid projectID length %v < 1 or > %v", len(ti.ProjectID), timestamp.MaxProjectIDLen)
	}

	return nil
}

func (ti TimeStampLastIndexStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice([]byte(ti.ProjectID), util.Uint64ToBytes(ti.Index))
}

func StateTimeStampLastIndexValue(st base.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(TimeStampLastIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.Index, nil
}

func IsStateTimeStampLastIndexKey(key string) bool {
	return strings.HasSuffix(key, StateKeyProjectLastIndexSuffix)
}

func StateKeyTimeStampLastIndex(addr base.Address, sid extensioncurrency.ContractID, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyTimeStampService(addr, sid), pid, StateKeyProjectLastIndexSuffix)
}

var (
	TimeStampItemStateValueHint = hint.MustNewHint("mitum-timestamp-item-state-value-v0.0.1")
	StateKeyTimeStampItemSuffix = ":timestampitem"
)

type TimeStampItemStateValue struct {
	hint.BaseHinter
	TimeStampItem timestamp.TimeStampItem
}

func NewTimeStampItemStateValue(item timestamp.TimeStampItem) TimeStampItemStateValue {
	return TimeStampItemStateValue{
		BaseHinter:    hint.NewBaseHinter(TimeStampItemStateValueHint),
		TimeStampItem: item,
	}
}

func (ts TimeStampItemStateValue) Hint() hint.Hint {
	return ts.BaseHinter.Hint()
}

func (ts TimeStampItemStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TimeStampItemStateValue")

	if err := ts.BaseHinter.IsValid(TimeStampItemStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ts.TimeStampItem.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ts TimeStampItemStateValue) HashBytes() []byte {
	return ts.TimeStampItem.Bytes()
}

func StateTimeStampItemValue(st base.State) (timestamp.TimeStampItem, error) {
	v := st.Value()
	if v == nil {
		return timestamp.TimeStampItem{}, util.ErrNotFound.Errorf("TimeStampItem not found in State")
	}

	ts, ok := v.(TimeStampItemStateValue)
	if !ok {
		return timestamp.TimeStampItem{}, errors.Errorf("invalid TimeStampItem value found, %T", v)
	}

	return ts.TimeStampItem, nil
}

func IsStateTimeStampItemKey(key string) bool {
	return strings.HasSuffix(key, StateKeyTimeStampItemSuffix)
}

func StateKeyTimeStampItem(addr base.Address, sid extensioncurrency.ContractID, pid string, index uint64) string {
	return fmt.Sprintf("%s:%s:s%s%s", StateKeyTimeStampService(addr, sid), pid, strconv.FormatUint(index, 10), StateKeyTimeStampItemSuffix)
}

func ParseStateKey(key string) (string, string, error) {
	parsedKey := strings.Split(key, ":")
	if parsedKey[0] != StateKeyTimeStampPrefix[:len(StateKeyTimeStampPrefix)-1] {
		return "", "", errors.Errorf("State Key not include TimeStampPrefix, %s", parsedKey)
	}
	if len(parsedKey) < 3 {
		return "", "", errors.Errorf("parsing State Key string failed, %s", parsedKey)
	} else {
		return parsedKey[1], parsedKey[2], nil
	}
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
		st = currency.NewBaseState(base.NilHeight, k, nil, nil, nil)
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
