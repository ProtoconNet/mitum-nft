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

type TestMintProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Mint, MintItem]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestMintProcessor(tp *test.TestProcessor) TestMintProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Mint, MintItem](tp)
	return TestMintProcessor{
		BaseTestOperationProcessorWithItem: &t,
	}
}

func (t *TestMintProcessor) Create() *TestMintProcessor {
	t.Opr, _ = NewMintProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestMintProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestMintProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestMintProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestMintProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestMintProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators types.Signers) *TestMintProcessor {
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

func (t *TestMintProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []types.Signer,
) *TestMintProcessor {
	sg := types.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[types.Signer](sg, target)

	return t
}

func (t *TestMintProcessor) SetSigners(
	signers []types.Signer, target []types.Signers,
) *TestMintProcessor {
	sg := types.NewSigners(signers)
	test.UpdateSlice[types.Signers](sg, target)

	return t
}

func (t *TestMintProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestMintProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestMintProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestMintProcessor {
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

func (t *TestMintProcessor) LoadOperation(fileName string,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestMintProcessor) Print(fileName string,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestMintProcessor) MakeItem(
	target test.Account, receiver test.Account, hash, uri string, creators types.Signers, currency ctypes.CurrencyID,
	targetItems []MintItem,
) *TestMintProcessor {
	item := NewMintItem(target.Address(), receiver.Address(), types.NFTHash(hash), types.URI(uri), creators, currency)
	test.UpdateSlice[MintItem](item, targetItems)

	return t
}

func (t *TestMintProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []MintItem,
) *TestMintProcessor {
	op, _ := NewMint(
		NewMintFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestMintProcessor) RunPreProcess() *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestMintProcessor) RunProcess() *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestMintProcessor) IsValid() *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestMintProcessor) Decode(fileName string) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
