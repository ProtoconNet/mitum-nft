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

type testCollectionPolicyUpdaterOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testCollectionPolicyUpdaterOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testCollectionPolicyUpdaterOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(CollectionPolicyUpdaterHinter, NewCollectionPolicyUpdaterProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testCollectionPolicyUpdaterOperations) newCollectionPolicyUpdater(sender base.Address, keys []key.Privatekey, symbol extensioncurrency.ContractID, policy CollectionPolicy, cid currency.CurrencyID) CollectionPolicyUpdater {
	token := util.UUID().Bytes()
	fact := NewCollectionPolicyUpdaterFact(token, sender, symbol, policy, cid)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	cpu, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	t.NoError(cpu.IsValid(nil))

	return cpu
}

func (t *testCollectionPolicyUpdaterOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, sender.Privs(), t.symbol, policy, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testCollectionPolicyUpdaterOperations) TestCollectionNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	sts = append(sts, sst...)

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, sender.Privs(), t.symbol, policy, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testCollectionPolicyUpdaterOperations) TestCollectionDeactivated() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(false, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, sender.Privs(), t.symbol, policy, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "deactivated collection")
}

func (t *testCollectionPolicyUpdaterOperations) TestSenderUnathorized() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	creator, cst := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, creator.Address)

	sts = append(sts, sst...)
	sts = append(sts, cst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, creator.Address, []base.Address{creator.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, sender.Privs(), t.symbol, policy, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not creator of collection design")
}

func (t *testCollectionPolicyUpdaterOperations) TestOperationWithFee() {
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
	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	fact := NewCollectionPolicyUpdaterFact(token, sender.Address, t.symbol, policy, t.cid)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	cpu, err := NewCollectionPolicyUpdater(fact, fs, "")
	t.NoError(err)

	err = opr.Process(cpu)
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

func (t *testCollectionPolicyUpdaterOperations) TestInSufficientBalanceWithFee() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(33), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	fee := currency.NewBig(34)
	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, sender.Privs(), t.symbol, policy, t.cid)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)
	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testCollectionPolicyUpdaterOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	_, dst0 := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst0...)

	_, dst1 := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, extensioncurrency.ContractID("ABC"), []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst1...)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token0 := util.UUID().Bytes()
	policy0 := NewCollectionPolicy("Collection0", 0, "", []base.Address{})
	fact0 := NewCollectionPolicyUpdaterFact(token0, sender.Address, t.symbol, policy0, t.cid)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	cpu0, err := NewCollectionPolicyUpdater(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(cpu0))

	token1 := util.UUID().Bytes()
	policy1 := NewCollectionPolicy("Collection1", 1, "", []base.Address{})
	fact1 := NewCollectionPolicyUpdaterFact(token1, sender.Address, extensioncurrency.ContractID("ABC"), policy1, t.cid)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	cpu1, err := NewCollectionPolicyUpdater(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(cpu1)

	t.Contains(err.Error(), "violates only one sender")
}

// func (t *testCollectionPolicyUpdaterOperations) TestSameCollection() {}

func (t *testCollectionPolicyUpdaterOperations) TestUnderThreshold() {
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

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender, pks, t.symbol, policy, t.cid)

	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testCollectionPolicyUpdaterOperations) TestUnknownKey() {
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

	policy := NewCollectionPolicy("Collection", 0, "", []base.Address{})
	cpu := t.newCollectionPolicyUpdater(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, t.symbol, policy, t.cid)

	err := opr.Process(cpu)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestCollectionPolicyUpdaterOperations(t *testing.T) {
	suite.Run(t, new(testCollectionPolicyUpdaterOperations))
}
