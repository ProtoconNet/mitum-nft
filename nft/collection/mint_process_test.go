package collection

import (
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util"
)

type testMintOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testMintOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testMintOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(MintHinter, NewMintProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testMintOperations) newMintItem(symbol extensioncurrency.ContractID, form MintForm, cid currency.CurrencyID) MintItem {
	return NewMintItem(symbol, form, cid)
}

func (t *testMintOperations) newMint(sender base.Address, keys []key.Privatekey, items []MintItem) Mint {
	token := util.UUID().Bytes()
	fact := NewMintFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	t.NoError(mint.IsValid(nil))

	return mint
}

func (t *testMintOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []MintItem{t.newMintItem(
		t.symbol,
		NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})),
		t.cid,
	)}
	mint := t.newMint(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testMintOperations) TestCollectionNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	sts = append(sts, sst...)

	items := []MintItem{t.newMintItem(
		t.symbol,
		NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})),
		t.cid,
	)}
	mint := t.newMint(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testMintOperations) TestCollectionDeactivated() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(false, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []MintItem{t.newMintItem(
		t.symbol,
		NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})),
		t.cid,
	)}
	mint := t.newMint(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "deactivated collection")
}

func (t *testMintOperations) TestUnauthorizedSender() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	owner, ost := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	sts = append(sts, sst...)
	sts = append(sts, ost...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, owner.Address, []base.Address{owner.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []MintItem{t.newMintItem(
		t.symbol,
		NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})),
		t.cid,
	)}
	mint := t.newMint(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "sender is not whitelisted")
}

func (t *testMintOperations) TestMaxCollectionIdx() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	sts = append(sts, sst...)
	sts = append(sts, pst)

	nfts := []nft.NFTID{}
	var i uint64 = 1
	for ; i < nft.MaxNFTIdx+1; i++ {
		nfts = append(nfts, nft.NewNFTID(t.symbol, i))
	}

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, nfts, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []MintItem{t.newMintItem(
		t.symbol,
		NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})),
		t.cid,
	)}
	mint := t.newMint(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "idx over max")
}

func (t *testMintOperations) TestMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []MintItem{
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/1", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/2", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
	}
	fact := NewMintFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	err = opr.Process(mint)
	t.NoError(err)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)

	var amst state.State
	var nftst0 state.State
	var nftst1 state.State
	var nboxst state.State
	var am currency.Amount
	var nf0 nft.NFT
	var nf1 nft.NFT
	var nbox NFTBox
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyNFT(nid0) {
			nftst0 = st.GetState()
			nf0, _ = StateNFTValue(nftst0)
		} else if st.Key() == StateKeyNFT(nid1) {
			nftst1 = st.GetState()
			nf1, _ = StateNFTValue(nftst1)
		} else if st.Key() == StateKeyNFTs(t.symbol) {
			nboxst = st.GetState()
			nbox, _ = StateNFTsValue(nboxst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee.MulInt64(2)), am.Big())
	t.Equal(fee.MulInt64(2), amst.(currency.AmountState).Fee())

	t.True(nf0.Owner().Equal(sender.Address))
	t.True(nf1.Owner().Equal(sender.Address))
	t.True(nbox.Exists(nid0))
	t.True(nbox.Exists(nid1))
}

func (t *testMintOperations) TestInsufficientMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(17)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []MintItem{
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/1", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/2", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
	}
	fact := NewMintFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	err = opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testMintOperations) TestInSufficientBalanceWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(34)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []MintItem{
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
	}
	fact := NewMintFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	mint, err := NewMint(fact, fs, "")
	t.NoError(err)

	err = opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testMintOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token0 := util.UUID().Bytes()
	items0 := []MintItem{
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/1", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
	}
	fact0 := NewMintFact(token0, sender.Address, items0)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	mint0, err := NewMint(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(mint0))

	token1 := util.UUID().Bytes()
	items1 := []MintItem{
		t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft/2", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid),
	}
	fact1 := NewMintFact(token1, sender.Address, items1)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	mint1, err := NewMint(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(mint1)

	t.Contains(err.Error(), "violates only one sender")
}

func (t *testMintOperations) TestUnderThreshold() {
	spk := key.NewBasePrivatekey()
	apk := key.NewBasePrivatekey()

	skey := t.newKey(spk.Publickey(), 50)
	akey := t.newKey(apk.Publickey(), 50)
	skeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{skey, akey}, 100)

	pks := []key.Privatekey{spk}
	sender, _ := currency.NewAddressFromKeys(skeys)

	// set sender state
	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)

	parent, _, pst := t.newContractAccount(true, true, sender)
	_, dst := t.newCollectionDesign(true, parent, sender, []base.Address{sender}, t.symbol, []nft.NFTID{}, []nft.NFTID{})

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		pst,
	)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []MintItem{t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid)}
	mint := t.newMint(sender, pks, items)

	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testMintOperations) TestUnknownKey() {
	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})

	sts := []state.State{}
	sts = append(sts, sst...)
	sts = append(sts, pst)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []MintItem{t.newMintItem(t.symbol, NewMintForm("", "https://localhost:5000/nft", nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{})), t.cid)}

	mint := t.newMint(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, items)

	err := opr.Process(mint)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestMintOperations(t *testing.T) {
	suite.Run(t, new(testMintOperations))
}
