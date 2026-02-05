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

type TestApproveProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Approve, ApproveItem]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestApproveProcessor(tp *test.TestProcessor) TestApproveProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Approve, ApproveItem](tp)
	return TestApproveProcessor{BaseTestOperationProcessorWithItem: &t}
}

func (t *TestApproveProcessor) Create() *TestApproveProcessor {
	t.Opr, _ = NewApproveProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestApproveProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestApproveProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestApproveProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestApproveProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestApproveProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []types.Signer,
) *TestApproveProcessor {
	sg := types.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[types.Signer](sg, target)

	return t
}

func (t *TestApproveProcessor) SetSigners(
	signers []types.Signer, target []types.Signers,
) *TestApproveProcessor {
	sg := types.NewSigners(signers)
	test.UpdateSlice[types.Signers](sg, target)

	return t
}

func (t *TestApproveProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators types.Signers) *TestApproveProcessor {
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

func (t *TestApproveProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestApproveProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	policy := types.NewCollectionPolicy(t.name, t.royalty, t.uri, whs)
	design := types.NewDesign(contract, sender, true, 0, policy)

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

	status.SetActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(status), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestApproveProcessor) LoadOperation(fileName string,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestApproveProcessor) Print(fileName string,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestApproveProcessor) MakeItem(
	target test.Account, approved test.Account, idx uint64, currency ctypes.CurrencyID, targetItems []ApproveItem,
) *TestApproveProcessor {
	item := NewApproveItem(target.Address(), approved.Address(), idx, currency)
	test.UpdateSlice[ApproveItem](item, targetItems)

	return t
}

func (t *TestApproveProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []ApproveItem,
) *TestApproveProcessor {
	op, _ := NewApprove(
		NewApproveFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestApproveProcessor) RunPreProcess() *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestApproveProcessor) RunProcess() *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestApproveProcessor) IsValid() *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestApproveProcessor) Decode(fileName string) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
