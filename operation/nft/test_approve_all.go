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

type TestDelegateProcessor struct {
	*test.BaseTestOperationProcessorWithItem[ApproveAll, ApproveAllItem]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestDelegateProcessor(tp *test.TestProcessor) TestDelegateProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[ApproveAll, ApproveAllItem](tp)
	return TestDelegateProcessor{BaseTestOperationProcessorWithItem: &t}
}

func (t *TestDelegateProcessor) Create() *TestDelegateProcessor {
	t.Opr, _ = NewDelegateProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestDelegateProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestDelegateProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestDelegateProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestDelegateProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestDelegateProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestDelegateProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestDelegateProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []types.Signer,
) *TestDelegateProcessor {
	sg := types.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[types.Signer](sg, target)

	return t
}

func (t *TestDelegateProcessor) SetSigners(
	signers []types.Signer, target []types.Signers,
) *TestDelegateProcessor {
	sg := types.NewSigners(signers)
	test.UpdateSlice[types.Signers](sg, target)

	return t
}

func (t *TestDelegateProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators types.Signers) *TestDelegateProcessor {
	cst, found, _ := t.MockGetter.Get(state.NFTStateKey(contract, state.LastIDXKey))
	if !found {
		panic("service not set")
	}

	nftID, _ := state.StateLastNFTIndexValue(cst)
	n := types.NewNFT(nftID, true, owner, types.NFTHash(nfthash), types.URI(uri), owner, creators)

	st := common.NewBaseState(base.Height(1), state.StateKeyNFT(contract, nftID), state.NewNFTStateValue(n), nil, []util.Hash{})
	t.SetState(st, true)
	st = common.NewBaseState(base.Height(1), state.NFTStateKey(contract, state.LastIDXKey), state.NewLastNFTIndexStateValue(nftID+1), nil, []util.Hash{})
	t.SetState(st, true)

	return t
}

func (t *TestDelegateProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestDelegateProcessor {
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

func (t *TestDelegateProcessor) LoadOperation(fileName string,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestDelegateProcessor) Print(fileName string,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestDelegateProcessor) MakeItem(
	target test.Account, operator test.Account, mode ApproveAllMode, currency ctypes.CurrencyID, targetItems []ApproveAllItem,
) *TestDelegateProcessor {
	item := NewApproveAllItem(target.Address(), operator.Address(), mode, currency)
	test.UpdateSlice[ApproveAllItem](item, targetItems)

	return t
}

func (t *TestDelegateProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []ApproveAllItem,
) *TestDelegateProcessor {
	op, _ := NewDelegate(
		NewApproveAllFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestDelegateProcessor) RunPreProcess() *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestDelegateProcessor) RunProcess() *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestDelegateProcessor) IsValid() *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestDelegateProcessor) Decode(fileName string) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
