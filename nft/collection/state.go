package collection

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var StateKeyAgentsSuffix = ":Agents"
var StateKeyCollectionPolicySuffix = ":Collection"

func StateKeyAgents(a base.Address) string {
	return fmt.Sprintf("%s%s", a.String(), StateKeyAgentsSuffix)
}

func IsStateAgentsKey(key string) bool {
	return strings.HasSuffix(key, StateKeyAgentsSuffix)
}

func StateAgentsValue(st state.State) (AgentBox, error) {
	v := st.Value()
	if v == nil {
		return AgentBox{}, util.NotFoundError.Errorf("agent box not found in State")
	}

	if s, ok := v.Interface().(AgentBox); !ok {
		return AgentBox{}, errors.Errorf("invalid agent box value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateAgentsValue(st state.State, b AgentBox) (state.State, error) {
	if uv, err := state.NewHintedValue(b); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func StateKeyCollectionPolicy(policy CollectionPolicy) string {
	return fmt.Sprintf("%s%s", policy.Symbol().String(), StateKeyCollectionPolicySuffix)
}

func IsStateCollectionPolicyKey(key string) bool {
	return strings.HasSuffix(key, StateKeyCollectionPolicySuffix)
}

func StateCollectionPolicyValue(st state.State) (CollectionPolicyBox, error) {
	v := st.Value()
	if v == nil {
		return CollectionPolicyBox{}, util.NotFoundError.Errorf("collection policy box not found in State")
	}

	if s, ok := v.Interface().(CollectionPolicyBox); !ok {
		return CollectionPolicyBox{}, errors.Errorf("invalid collection policy box value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateCollectionPolicyValue(st state.State, b CollectionPolicyBox) (state.State, error) {
	if uv, err := state.NewHintedValue(b); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
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
