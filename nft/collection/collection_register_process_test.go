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

type testCollectionRegisterOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testCollectionRegisterOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testCollectionRegisterOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(CollectionRegisterHinter, NewCollectionRegisterProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testCollectionRegisterOperations) newCollectionRegister(sender base.Address, keys []key.Privatekey, form CollectionRegisterForm, cid currency.CurrencyID) CollectionRegister {
	token := util.UUID().Bytes()
	fact := NewCollectionRegisterFact(token, sender, form, cid)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	cr, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	t.NoError(cr.IsValid(nil))

	return cr
}

func (t *testCollectionRegisterOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	sts = append(sts, pst)

	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{sender.Address})
	cr := t.newCollectionRegister(sender.Address, sender.Privs(), form, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testCollectionRegisterOperations) TestParentNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, _ := t.newContractAccount(false, true, sender.Address)
	sts = append(sts, sst...)

	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{sender.Address})
	cr := t.newCollectionRegister(sender.Address, sender.Privs(), form, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testCollectionRegisterOperations) TestParentDeactivated() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, false, sender.Address)
	sts = append(sts, sst...)
	sts = append(sts, pst)

	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{sender.Address})
	cr := t.newCollectionRegister(sender.Address, sender.Privs(), form, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "deactivated contract account")
}

func (t *testCollectionRegisterOperations) TestOperationWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{})
	fact := NewCollectionRegisterFact(token, sender.Address, form, t.cid)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	cr, err := NewCollectionRegister(fact, fs, "")
	t.NoError(err)

	err = opr.Process(cr)
	t.NoError(err)

	var amst state.State
	var am currency.Amount
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())
}

func (t *testCollectionRegisterOperations) TestInSufficientBalanceWithFee() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(33), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	fee := currency.NewBig(34)
	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{})
	cr := t.newCollectionRegister(sender.Address, sender.Privs(), form, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testCollectionRegisterOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token0 := util.UUID().Bytes()
	form0 := NewCollectionRegisterForm(parent, t.symbol, "Collection0", 0, "", []base.Address{})
	fact0 := NewCollectionRegisterFact(token0, sender.Address, form0, t.cid)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	cpu0, err := NewCollectionRegister(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(cpu0))

	token1 := util.UUID().Bytes()
	form1 := NewCollectionRegisterForm(parent, extensioncurrency.ContractID("ABC"), "Collection1", 1, "", []base.Address{})
	fact1 := NewCollectionRegisterFact(token1, sender.Address, form1, t.cid)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	cpu1, err := NewCollectionRegister(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(cpu1)

	t.Contains(err.Error(), "violates only one sender")
}

func (t *testCollectionRegisterOperations) TestUnderThreshold() {
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

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		pst,
	)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{})
	cr := t.newCollectionRegister(sender, pks, form, t.cid)

	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testCollectionRegisterOperations) TestUnknownKey() {
	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts := []state.State{}
	sts = append(sts, sst...)
	sts = append(sts, pst)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	form := NewCollectionRegisterForm(parent, t.symbol, "Collection", 0, "", []base.Address{})
	cr := t.newCollectionRegister(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, form, t.cid)

	err := opr.Process(cr)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestCollectionRegisterOperations(t *testing.T) {
	suite.Run(t, new(testCollectionRegisterOperations))
}
