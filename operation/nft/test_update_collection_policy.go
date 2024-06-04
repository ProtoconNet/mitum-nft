package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	statenft "github.com/ProtoconNet/mitum-nft/state"
	nfttypes "github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TestUpdateCollectionPolicyProcessor struct {
	*test.BaseTestOperationProcessorNoItem[UpdateCollectionPolicy]
	name    nfttypes.CollectionName
	royalty nfttypes.PaymentParameter
	uri     nfttypes.URI
}

func NewTestUpdateCollectionPolicyProcessor(tp *test.TestProcessor) TestUpdateCollectionPolicyProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[UpdateCollectionPolicy](tp)
	return TestUpdateCollectionPolicyProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestUpdateCollectionPolicyProcessor) Create() *TestUpdateCollectionPolicyProcessor {
	t.Opr, _ = NewUpdateCollectionPolicyProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestUpdateCollectionPolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
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
	t.name = nfttypes.CollectionName(name)
	t.royalty = nfttypes.PaymentParameter(royalty)
	t.uri = nfttypes.URI(uri)

	return t
}

func (t *TestUpdateCollectionPolicyProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestUpdateCollectionPolicyProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	policy := nfttypes.NewCollectionPolicy(t.name, t.royalty, t.uri, whs)
	design := nfttypes.NewDesign(contract, sender, true, policy)

	st := common.NewBaseState(base.Height(1), statenft.NFTStateKey(design.Parent(), statenft.CollectionKey), statenft.NewCollectionStateValue(design), nil, []util.Hash{})
	t.SetState(st, true)
	st = common.NewBaseState(base.Height(1), statenft.NFTStateKey(design.Parent(), statenft.LastIDXKey), statenft.NewLastNFTIndexStateValue(0), nil, []util.Hash{})
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
	sender base.Address, privatekey base.Privatekey, contract base.Address, whitelist []test.Account, currency types.CurrencyID,
) *TestUpdateCollectionPolicyProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	op, _ := NewUpdateCollectionPolicy(
		NewUpdateCollectionPolicyFact(
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
