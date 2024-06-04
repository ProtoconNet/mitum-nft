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

type TestDelegateProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Delegate, DelegateItem]
	name    nfttypes.CollectionName
	royalty nfttypes.PaymentParameter
	uri     nfttypes.URI
}

func NewTestDelegateProcessor(tp *test.TestProcessor) TestDelegateProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Delegate, DelegateItem](tp)
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
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestDelegateProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestDelegateProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestDelegateProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestDelegateProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestDelegateProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestDelegateProcessor {
	t.name = nfttypes.CollectionName(name)
	t.royalty = nfttypes.PaymentParameter(royalty)
	t.uri = nfttypes.URI(uri)

	return t
}

func (t *TestDelegateProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []nfttypes.Signer,
) *TestDelegateProcessor {
	sg := nfttypes.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[nfttypes.Signer](sg, target)

	return t
}

func (t *TestDelegateProcessor) SetSigners(
	signers []nfttypes.Signer, target []nfttypes.Signers,
) *TestDelegateProcessor {
	sg := nfttypes.NewSigners(signers)
	test.UpdateSlice[nfttypes.Signers](sg, target)

	return t
}

func (t *TestDelegateProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators nfttypes.Signers) *TestDelegateProcessor {
	cst, found, _ := t.MockGetter.Get(statenft.NFTStateKey(contract, statenft.LastIDXKey))
	if !found {
		panic("service not set")
	}

	nftID, _ := statenft.StateLastNFTIndexValue(cst)
	n := nfttypes.NewNFT(nftID, true, owner, nfttypes.NFTHash(nfthash), nfttypes.URI(uri), owner, creators)

	st := common.NewBaseState(base.Height(1), statenft.StateKeyNFT(contract, nftID), statenft.NewNFTStateValue(n), nil, []util.Hash{})
	t.SetState(st, true)
	st = common.NewBaseState(base.Height(1), statenft.NFTStateKey(contract, statenft.LastIDXKey), statenft.NewLastNFTIndexStateValue(nftID+1), nil, []util.Hash{})
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
	target test.Account, operator test.Account, mode DelegateMode, currency types.CurrencyID, targetItems []DelegateItem,
) *TestDelegateProcessor {
	item := NewDelegateItem(target.Address(), operator.Address(), mode, currency)
	test.UpdateSlice[DelegateItem](item, targetItems)

	return t
}

func (t *TestDelegateProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []DelegateItem,
) *TestDelegateProcessor {
	op, _ := NewDelegate(
		NewDelegateFact(
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
