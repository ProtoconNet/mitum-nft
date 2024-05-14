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
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TestApproveProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Approve, ApproveItem]
	name    nfttypes.CollectionName
	royalty nfttypes.PaymentParameter
	uri     nfttypes.URI
}

func NewTestApproveProcessor(encs *encoder.Encoders) TestApproveProcessor {
	t := test.NewBaseTestOperationProcessorWithItem[Approve, ApproveItem](encs)
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
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestApproveProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestApproveProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestApproveProcessor {
	t.name = nfttypes.CollectionName(name)
	t.royalty = nfttypes.PaymentParameter(royalty)
	t.uri = nfttypes.URI(uri)

	return t
}

func (t *TestApproveProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []nfttypes.Signer,
) *TestApproveProcessor {
	sg := nfttypes.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[nfttypes.Signer](sg, target)

	return t
}

func (t *TestApproveProcessor) SetSigners(
	signers []nfttypes.Signer, target []nfttypes.Signers,
) *TestApproveProcessor {
	sg := nfttypes.NewSigners(signers)
	test.UpdateSlice[nfttypes.Signers](sg, target)

	return t
}

func (t *TestApproveProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators nfttypes.Signers) *TestApproveProcessor {
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

func (t *TestApproveProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestApproveProcessor {
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
	target test.Account, approved test.Account, idx uint64, currency types.CurrencyID, targetItems []ApproveItem,
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
