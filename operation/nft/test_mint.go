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

type TestMintProcessor struct {
	*test.BaseTestOperationProcessorWithItem[Mint, MintItem]
	name    nfttypes.CollectionName
	royalty nfttypes.PaymentParameter
	uri     nfttypes.URI
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
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestMintProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetAmount(am, cid, target)

	return t
}

func (t *TestMintProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestMintProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestMintProcessor {
	t.BaseTestOperationProcessorWithItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestMintProcessor) SetNFT(contract, owner base.Address, nfthash, uri string, creators nfttypes.Signers) *TestMintProcessor {
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

func (t *TestMintProcessor) SetSigner(
	signer test.Account, share uint, signed bool, target []nfttypes.Signer,
) *TestMintProcessor {
	sg := nfttypes.NewSigner(signer.Address(), share, signed)
	test.UpdateSlice[nfttypes.Signer](sg, target)

	return t
}

func (t *TestMintProcessor) SetSigners(
	signers []nfttypes.Signer, target []nfttypes.Signers,
) *TestMintProcessor {
	sg := nfttypes.NewSigners(signers)
	test.UpdateSlice[nfttypes.Signers](sg, target)

	return t
}

func (t *TestMintProcessor) SetDesign(
	name string,
	royalty uint,
	uri string,
) *TestMintProcessor {
	t.name = nfttypes.CollectionName(name)
	t.royalty = nfttypes.PaymentParameter(royalty)
	t.uri = nfttypes.URI(uri)

	return t
}

func (t *TestMintProcessor) SetService(
	sender, contract base.Address, whitelist []test.Account,
) *TestMintProcessor {
	var whs []base.Address
	for _, wh := range whitelist {
		whs = append(whs, wh.Address())
	}

	policy := nfttypes.NewCollectionPolicy(t.name, t.royalty, t.uri, whs)
	design := nfttypes.NewDesign(contract, sender, true, policy)

	st := common.NewBaseState(base.Height(1), statenft.NFTStateKey(design.Contract(), statenft.CollectionKey), statenft.NewCollectionStateValue(design), nil, []util.Hash{})
	t.SetState(st, true)
	st = common.NewBaseState(base.Height(1), statenft.NFTStateKey(design.Contract(), statenft.LastIDXKey), statenft.NewLastNFTIndexStateValue(0), nil, []util.Hash{})
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
	target test.Account, receiver test.Account, hash, uri string, creators nfttypes.Signers, currency types.CurrencyID,
	targetItems []MintItem,
) *TestMintProcessor {
	item := NewMintItem(target.Address(), receiver.Address(), nfttypes.NFTHash(hash), nfttypes.URI(uri), creators, currency)
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
