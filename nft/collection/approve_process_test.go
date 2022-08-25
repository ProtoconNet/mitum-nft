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

type testApproveOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testApproveOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testApproveOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(ApproveHinter, NewApproveProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testApproveOperations) newApproveItem(approved base.Address, nid nft.NFTID, cid currency.CurrencyID) ApproveItem {
	return NewApproveItem(approved, nid, cid)
}

func (t *testApproveOperations) newApprove(sender base.Address, keys []key.Privatekey, items []ApproveItem) Approve {
	token := util.UUID().Bytes()
	fact := NewApproveFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	t.NoError(approve.IsValid(nil))

	return approve
}

func (t *testApproveOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, ast...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, sender.Address, "", "https://localhost:5000/nft", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testApproveOperations) TestApprovedNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, _ := t.newAccount(false, nil)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, sender.Address, "", "https://localhost:5000/nft", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testApproveOperations) TestNFTNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, ast...)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testApproveOperations) TestNFTBurned() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, ast...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, false, sender.Address, "", "https://localhost:5000/nft", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{nid})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "burned nft")
}

func (t *testApproveOperations) TestAgentApprove() {
	var sts = []state.State{}

	agentBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	owner, sst := t.newAccount(true, nil)
	agent, agst := t.newAccount(true, []currency.Amount{agentBalance})
	parent, _, pst := t.newContractAccount(true, true, owner.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, agst...)
	sts = append(sts, ast...)

	boxst := t.newStateAgent(owner.Address, t.symbol, []base.Address{agent.Address})
	sts = append(sts, boxst)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, owner.Address, "", "https://localhost:5000/nft", owner.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, owner.Address, []base.Address{owner.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(agent.Address, agent.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(owner.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)
	t.NoError(err)

	var amst state.State
	var nftst state.State
	var am currency.Amount
	var nf nft.NFT
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(agent.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyNFT(nid) {
			nftst = st.GetState()
			nf, _ = StateNFTValue(nftst)
		}
	}

	t.Equal(agentBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())

	t.True(nf.Approved().Equal(approved.Address))
}

func (t *testApproveOperations) TestUnauthorizedSender() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	owner, ost := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, ost...)
	sts = append(sts, ast...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, owner.Address, "", "https://localhost:5000/nft", owner.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, owner.Address, []base.Address{owner.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}
	approve := t.newApprove(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unauthorized sender")
}

func (t *testApproveOperations) TestMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	approved0, ast0 := t.newAccount(true, nil)
	approved1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(nid0, true, sender.Address, "", "https://localhost:5000/nft/1", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	n1 := nft.NewNFT(nid1, true, sender.Address, "", "https://localhost:5000/nft/2", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst0 := t.newStateNFT(n0)
	nst1 := t.newStateNFT(n1)

	sts = append(sts, nst0, nst1)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid0, nid1}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []ApproveItem{
		t.newApproveItem(approved0.Address, nid0, t.cid),
		t.newApproveItem(approved1.Address, nid1, t.cid),
	}
	fact := NewApproveFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = opr.Process(approve)
	t.NoError(err)

	var amst state.State
	var nftst0 state.State
	var nftst1 state.State
	var am currency.Amount
	var nf0 nft.NFT
	var nf1 nft.NFT
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
		}
	}

	t.Equal(senderBalance.Big().Sub(fee.MulInt64(2)), am.Big())
	t.Equal(fee.MulInt64(2), amst.(currency.AmountState).Fee())

	t.True(nf0.Approved().Equal(approved0.Address))
	t.True(nf1.Approved().Equal(approved1.Address))
}

func (t *testApproveOperations) TestInsufficientMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	approved0, ast0 := t.newAccount(true, nil)
	approved1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(nid0, true, sender.Address, "", "https://localhost:5000/nft/1", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	n1 := nft.NewNFT(nid1, true, sender.Address, "", "https://localhost:5000/nft/2", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst0 := t.newStateNFT(n0)
	nst1 := t.newStateNFT(n1)

	sts = append(sts, nst0, nst1)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid0, nid1}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(17)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []ApproveItem{
		t.newApproveItem(approved0.Address, nid0, t.cid),
		t.newApproveItem(approved1.Address, nid1, t.cid),
	}
	fact := NewApproveFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testApproveOperations) TestInSufficientBalanceWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved, ast := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast...)
	sts = append(sts, pst)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, sender.Address, "", "https://localhost:5000/nft/1", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst := t.newStateNFT(n)

	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(34)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []ApproveItem{
		t.newApproveItem(approved.Address, nid, t.cid),
	}
	fact := NewApproveFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	approve, err := NewApprove(fact, fs, "")
	t.NoError(err)

	err = opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testApproveOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	approved0, ast0 := t.newAccount(true, nil)
	approved1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(nid0, true, sender.Address, "", "https://localhost:5000/nft/1", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))
	n1 := nft.NewNFT(nid1, true, sender.Address, "", "https://localhost:5000/nft/2", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst0 := t.newStateNFT(n0)
	nst1 := t.newStateNFT(n1)

	sts = append(sts, nst0, nst1)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid0, nid1}, []nft.NFTID{})
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token0 := util.UUID().Bytes()
	items0 := []ApproveItem{
		t.newApproveItem(approved0.Address, nid0, t.cid),
	}
	fact0 := NewApproveFact(token0, sender.Address, items0)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	approve0, err := NewApprove(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(approve0))

	token1 := util.UUID().Bytes()
	items1 := []ApproveItem{
		t.newApproveItem(approved1.Address, nid1, t.cid),
	}
	fact1 := NewApproveFact(token1, sender.Address, items1)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	approve1, err := NewApprove(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(approve1)

	t.Contains(err.Error(), "violates only one sender")
}

// func (t *testApproveOperations) TestSameNFTID() { }

func (t *testApproveOperations) TestUnderThreshold() {
	spk := key.NewBasePrivatekey()
	apk := key.NewBasePrivatekey()

	skey := t.newKey(spk.Publickey(), 50)
	akey := t.newKey(apk.Publickey(), 50)
	skeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{skey, akey}, 100)
	akeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{akey}, 50)

	pks := []key.Privatekey{spk}
	sender, _ := currency.NewAddressFromKeys(skeys)
	approved, _ := currency.NewAddressFromKeys(akeys)

	// set sender state
	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)

	parent, _, pst := t.newContractAccount(true, true, sender)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, sender, "", "https://localhost:5000/nft/1", sender, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst := t.newStateNFT(n)
	_, dst := t.newCollectionDesign(true, parent, sender, []base.Address{sender}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		t.newStateKeys(approved, akeys),
		pst,
		nst,
	)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []ApproveItem{t.newApproveItem(approved, nid, t.cid)}
	approve := t.newApprove(sender, pks, items)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testApproveOperations) TestUnknownKey() {
	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})
	approved, ast := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})

	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(nid, true, sender.Address, "", "https://localhost:5000/nft/1", sender.Address, nft.NewSigners(0, []nft.Signer{}), nft.NewSigners(0, []nft.Signer{}))

	nst := t.newStateNFT(n)
	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})

	sts := []state.State{}
	sts = append(sts, sst...)
	sts = append(sts, ast...)
	sts = append(sts, pst)
	sts = append(sts, dst...)
	sts = append(sts, nst)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []ApproveItem{t.newApproveItem(approved.Address, nid, t.cid)}

	approve := t.newApprove(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, items)

	err := opr.Process(approve)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestApproveOperations(t *testing.T) {
	suite.Run(t, new(testApproveOperations))
}
