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

type TestTransferProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Transfer, TransferItem]
	name    types.CollectionName
	royalty types.PaymentParameter
	uri     types.URI
}

func NewTestTransferProcessor(tp *test.TestProcessor) TestTransferProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Transfer, TransferItem](tp)
	return TestTransferProcessor{
		BaseTestOperationProcessorWithItem: &t,
	}
}

func (t *TestTransferProcessor) Create() *TestTransferProcessor {
	t.Opr, _ = NewTransferProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestTransferProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []ctypes.CurrencyID, instate bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestTransferProcessor) SetAmount(
	am int64, cid ctypes.CurrencyID, target []ctypes.Amount,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestTransferProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferProcessor) SetAccount(
	priv string, amount int64, cid ctypes.CurrencyID, target []test.Account, inState bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferProcessor) SetNFT(
	contract, owner base.Address, nfthash, uri string, creators types.Signers,
) *TestTransferProcessor {
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

func (t *TestTransferProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []types.Signer,
) *TestTransferProcessor {
	sg := types.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[types.Signer](sg, target)

	return t
}

func (t *TestTransferProcessor) SetSigners(
	signers []types.Signer, target []types.Signers,
) *TestTransferProcessor {
	sg := types.NewSigners(signers)
	test.UpdateSlice[types.Signers](sg, target)

	return t
}

func (t *TestTransferProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestTransferProcessor {
	t.name = types.CollectionName(name)
	t.royalty = types.PaymentParameter(royalty)
	t.uri = types.URI(uri)

	return t
}

func (t *TestTransferProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestTransferProcessor {
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

func (t *TestTransferProcessor) LoadOperation(fileName string,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.LoadOperation(fileName)

	return t
}

func (t *TestTransferProcessor) Print(fileName string,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.Print(fileName)

	return t
}

func (t *TestTransferProcessor) MakeItem(
	target test.Account, receiver test.Account, nftID uint64, currency ctypes.CurrencyID,
	targetItems []TransferItem,
) *TestTransferProcessor {
	item := NewTransferItem(target.Address(), receiver.Address(), nftID, currency)
	test.UpdateSlice[TransferItem](item, targetItems)

	return t
}

func (t *TestTransferProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, items []TransferItem,
) *TestTransferProcessor {
	op, _ := NewTransfer(
		NewTransferFact(
			[]byte("token"),
			sender,
			items,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestTransferProcessor) RunPreProcess() *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.RunPreProcess()

	return t
}

func (t *TestTransferProcessor) RunProcess() *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.RunProcess()

	return t
}

func (t *TestTransferProcessor) IsValid() *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.IsValid()

	return t
}

func (t *TestTransferProcessor) Decode(fileName string) *TestTransferProcessor {
	t.BaseTestOperationProcessorWithItem.Decode(fileName)

	return t
}
