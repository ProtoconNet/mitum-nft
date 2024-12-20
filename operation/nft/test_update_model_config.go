package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TestUpdateCollectionPolicyProcessor struct {
	*test.BaseTestOperationProcessorNoItem[UpdateModelConfig]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestUpdateCollectionPolicyProcessor(tp *test.TestProcessor) TestUpdateCollectionPolicyProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[UpdateModelConfig](tp)
	return TestUpdateCollectionPolicyProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestUpdateCollectionPolicyProcessor) Create() *TestUpdateCollectionPolicyProcessor {
	t.Opr, _ = NewUpdateModelConfigProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) LoadOperation(fileName string,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) Print(fileName string,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestUpdateCollectionPolicyProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestUpdateCollectionPolicyProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	policy := types.NewCollectionPolicy(t.name, t.royalty, t.uri, whs)
	design := types.NewDesign(contract, sender, true, policy)

	st := common.NewBaseState(base.Height(1), state.NFTStateKey(design.Contract(), state.CollectionKey), state.NewCollectionStateValue(design), nil, []util.Hash{})
	t.SetState(st, true)
	st = common.NewBaseState(base.Height(1), state.NFTStateKey(design.Contract(), state.LastIDXKey), state.NewLastNFTIndexStateValue(0), nil, []util.Hash{})
	t.SetState(st, true)

	cst, found, _ := t.MockGetter.Get(extension.StateKeyContractAccount(contract))
	if !found {
		panic("contract account not set")
	}
	status, err := extension.StateContractAccountValue(cst)
	if err != nil {
		panic(err)
	}

	nstatus := status.SetIsActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(nstatus), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, whitelist []test.Account, currency ctypes.CurrencyID,
) *TestUpdateCollectionPolicyProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	op, _ := NewUpdateModelConfig(
		NewUpdateModelConfigFact(
			[]byte("token"),
			sender,
			contract,
			t.name,
			t.royalty,
			t.uri,
			whs,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) RunPreProcess() *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) RunProcess() *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) IsValid() *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) Decode(fileName string) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
