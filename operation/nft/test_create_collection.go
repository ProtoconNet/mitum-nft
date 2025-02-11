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

type TestCreateCollectionProcessor struct {
	*test.BaseTestOperationProcessorNoItem[RegisterModel]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestCreateCollectionProcessor(tp *test.TestProcessor) TestCreateCollectionProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[RegisterModel](tp)
	return TestCreateCollectionProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestCreateCollectionProcessor) Create() *TestCreateCollectionProcessor {
	t.Opr, _ = NewRegisterModelProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestCreateCollectionProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestCreateCollectionProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestCreateCollectionProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateCollectionProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateCollectionProcessor) LoadOperation(fileName string,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestCreateCollectionProcessor) Print(fileName string,
) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestCreateCollectionProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestCreateCollectionProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestCreateCollectionProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestCreateCollectionProcessor {
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

	nstatus := status.SetActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(nstatus), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestCreateCollectionProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, whitelist []test.Account, currency ctypes.CurrencyID,
) *TestCreateCollectionProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	op, _ := NewRegisterModel(
		NewRegisterModelFact(
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

func (t *TestCreateCollectionProcessor) RunPreProcess() *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestCreateCollectionProcessor) RunProcess() *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestCreateCollectionProcessor) IsValid() *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestCreateCollectionProcessor) Decode(fileName string) *TestCreateCollectionProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
