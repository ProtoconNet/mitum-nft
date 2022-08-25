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

type testDelegateOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testDelegateOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testDelegateOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(DelegateHinter, NewDelegateProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testDelegateOperations) newDelegateItem(symbol extensioncurrency.ContractID, agent base.Address, mode DelegateMode, cid currency.CurrencyID) DelegateItem {
	return NewDelegateItem(symbol, agent, mode, cid)
}

func (t *testDelegateOperations) newDelegate(sender base.Address, keys []key.Privatekey, items []DelegateItem) Delegate {
	token := util.UUID().Bytes()
	fact := NewDelegateFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	t.NoError(delegate.IsValid(nil))

	return delegate
}

func (t *testDelegateOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	agent, ast := t.newAccount(true, nil)

	sts = append(sts, pst)
	sts = append(sts, ast...)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testDelegateOperations) TestAgentNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	agent, _ := t.newAccount(false, nil)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testDelegateOperations) TestCollectionNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	agent, _ := t.newAccount(true, nil)

	sts = append(sts, sst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testDelegateOperations) TestCollectionDeactivated() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	agent, ast := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, pst)
	sts = append(sts, ast...)

	_, dst := t.newCollectionDesign(false, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "deactivated collection")
}

func (t *testDelegateOperations) TestDelegateAllow() {
	var sts = []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	agent, ast := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, ast...)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(delegate)
	t.NoError(err)

	var amst state.State
	var am currency.Amount
	var boxst state.State
	var agbox AgentBox
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyAgents(sender.Address, t.symbol) {
			boxst = st.GetState()
			agbox, _ = StateAgentsValue(boxst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())

	t.True(agbox.Exists(agent.Address))
}

func (t *testDelegateOperations) TestDelegateCancel() {
	var sts = []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	agent, ast := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, ast...)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	agst := t.newStateAgent(sender.Address, t.symbol, []base.Address{agent.Address})
	sts = append(sts, agst)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateCancel, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(delegate)
	t.NoError(err)

	var amst state.State
	var am currency.Amount
	var boxst state.State
	var agbox AgentBox
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyAgents(sender.Address, t.symbol) {
			boxst = st.GetState()
			agbox, _ = StateAgentsValue(boxst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())

	t.True(!agbox.Exists(agent.Address))
}

func (t *testDelegateOperations) TestDelegateCancelNotAllowed() {
	var sts = []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	agent, ast := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, ast...)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateCancel, t.cid)}
	delegate := t.newDelegate(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not found in agent box")
}

func (t *testDelegateOperations) TestMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	agent0, ast0 := t.newAccount(true, nil)
	agent1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
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
	items := []DelegateItem{
		t.newDelegateItem(t.symbol, agent0.Address, DelegateAllow, t.cid),
		t.newDelegateItem(t.symbol, agent1.Address, DelegateAllow, t.cid),
	}
	fact := NewDelegateFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = opr.Process(delegate)
	t.NoError(err)

	var amst state.State
	var agboxst state.State
	var am currency.Amount
	var agbox AgentBox
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyAgents(sender.Address, t.symbol) {
			agboxst = st.GetState()
			agbox, _ = StateAgentsValue(agboxst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee.MulInt64(2)), am.Big())
	t.Equal(fee.MulInt64(2), amst.(currency.AmountState).Fee())

	t.True(agbox.Exists(agent0.Address))
	t.True(agbox.Exists(agent1.Address))
}

func (t *testDelegateOperations) TestInsufficientMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	agent0, ast0 := t.newAccount(true, nil)
	agent1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
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
	items := []DelegateItem{
		t.newDelegateItem(t.symbol, agent0.Address, DelegateAllow, t.cid),
		t.newDelegateItem(t.symbol, agent1.Address, DelegateAllow, t.cid),
	}
	fact := NewDelegateFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testDelegateOperations) TestInSufficientBalanceWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	agent, ast := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast...)
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
	items := []DelegateItem{
		t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid),
	}
	fact := NewDelegateFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	delegate, err := NewDelegate(fact, fs, "")
	t.NoError(err)

	err = opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testDelegateOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	agent0, ast0 := t.newAccount(true, nil)
	agent1, ast1 := t.newAccount(true, nil)

	sts = append(sts, sst...)
	sts = append(sts, ast0...)
	sts = append(sts, ast1...)
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
	items0 := []DelegateItem{
		t.newDelegateItem(t.symbol, agent0.Address, DelegateAllow, t.cid),
	}
	fact0 := NewDelegateFact(token0, sender.Address, items0)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	delegate0, err := NewDelegate(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(delegate0))

	token1 := util.UUID().Bytes()
	items1 := []DelegateItem{
		t.newDelegateItem(t.symbol, agent1.Address, DelegateAllow, t.cid),
	}
	fact1 := NewDelegateFact(token1, sender.Address, items1)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	delegate1, err := NewDelegate(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(delegate1)

	t.Contains(err.Error(), "violates only one sender")
}

func (t *testDelegateOperations) TestUnderThreshold() {
	spk := key.NewBasePrivatekey()
	apk := key.NewBasePrivatekey()

	skey := t.newKey(spk.Publickey(), 50)
	akey := t.newKey(apk.Publickey(), 50)
	skeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{skey, akey}, 100)
	akeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{akey}, 50)

	pks := []key.Privatekey{spk}
	sender, _ := currency.NewAddressFromKeys(skeys)
	agent, _ := currency.NewAddressFromKeys(akeys)

	// set sender state
	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)

	parent, _, pst := t.newContractAccount(true, true, sender)
	_, dst := t.newCollectionDesign(true, parent, sender, []base.Address{sender}, t.symbol, []nft.NFTID{}, []nft.NFTID{})

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		t.newStateKeys(agent, akeys),
		pst,
	)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent, DelegateAllow, t.cid)}
	delegate := t.newDelegate(sender, pks, items)

	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testDelegateOperations) TestUnknownKey() {
	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})
	agent, ast := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})

	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})

	sts := []state.State{}
	sts = append(sts, sst...)
	sts = append(sts, ast...)
	sts = append(sts, pst)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []DelegateItem{t.newDelegateItem(t.symbol, agent.Address, DelegateAllow, t.cid)}

	delegate := t.newDelegate(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, items)

	err := opr.Process(delegate)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestDelegateOperations(t *testing.T) {
	suite.Run(t, new(testDelegateOperations))
}
