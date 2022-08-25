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

type testSignOperations struct {
	baseTestOperationProcessor
	cid    currency.CurrencyID
	symbol extensioncurrency.ContractID
}

func (t *testSignOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.symbol = extensioncurrency.ContractID("SCOLLECT")
}

func (t *testSignOperations) processor(cp *extensioncurrency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(SignHinter, NewSignProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testSignOperations) newSignItem(q Qualification, nid nft.NFTID, cid currency.CurrencyID) SignItem {
	return NewSignItem(q, nid, cid)
}

func (t *testSignOperations) newSign(sender base.Address, keys []key.Privatekey, items []SignItem) Sign {
	token := util.UUID().Bytes()
	fact := NewSignFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range keys {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	signOp, err := NewSign(fact, fs, "")
	t.NoError(err)

	t.NoError(signOp.IsValid(nil))

	return signOp
}

func (t *testSignOperations) TestSenderNotExist() {
	var sts = []state.State{}

	sender, _ := t.newAccount(false, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testSignOperations) TestNFTNotExist() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testSignOperations) TestAlreadySigned() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, true)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{nid})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "this signer has already signed nft")
}

func (t *testSignOperations) TestUnauthorizedSender() {
	var sts = []state.State{}

	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1000), t.cid)})
	creator, sgst := t.newAccount(true, nil)
	parent, _, pst := t.newContractAccount(true, true, creator.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)
	sts = append(sts, sgst...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		creator.Address,
		"",
		"https://localhost:5000/nft",
		creator.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(creator.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, creator.Address, []base.Address{creator.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(creator.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not signer of nft")
}

func (t *testSignOperations) TestSignCreator() {
	var sts = []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{nid})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)
	t.NoError(err)

	var amst state.State
	var nftst state.State
	var am currency.Amount
	var nf nft.NFT
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyNFT(nid) {
			nftst = st.GetState()
			nf, _ = StateNFTValue(nftst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())

	t.True(nf.Creators().IsSignedByAddress(sender.Address))
}

func (t *testSignOperations) TestSignCopyrighter() {
	var sts = []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(1000), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, pst)
	sts = append(sts, sst...)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{}),
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
	)
	nst := t.newStateNFT(n)
	sts = append(sts, nst)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{}, []nft.NFTID{nid})
	sts = append(sts, dst...)

	items := []SignItem{t.newSignItem(CopyrighterQualification, nid, t.cid)}
	signOp := t.newSign(sender.Address, sender.Privs(), items)

	pool, _ := t.statepool(sts)

	fee := currency.NewBig(2)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, fee, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	err := opr.Process(signOp)
	t.NoError(err)

	var amst state.State
	var nftst state.State
	var am currency.Amount
	var nf nft.NFT
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sender.Address, t.cid) {
			amst = st.GetState()
			am, _ = currency.StateBalanceValue(amst)
		} else if st.Key() == StateKeyNFT(nid) {
			nftst = st.GetState()
			nf, _ = StateNFTValue(nftst)
		}
	}

	t.Equal(senderBalance.Big().Sub(fee), am.Big())
	t.Equal(fee, amst.(currency.AmountState).Fee())

	t.True(nf.Copyrighters().IsSignedByAddress(sender.Address))
}

func (t *testSignOperations) TestMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(
		nid0,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/1",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	n1 := nft.NewNFT(
		nid1,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/2",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)

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
	items := []SignItem{
		t.newSignItem(CreatorQualification, nid0, t.cid),
		t.newSignItem(CreatorQualification, nid1, t.cid),
	}
	fact := NewSignFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	signOp, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = opr.Process(signOp)
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

	t.True(nf0.Creators().IsSignedByAddress(sender.Address))
	t.True(nf1.Creators().IsSignedByAddress(sender.Address))
}

func (t *testSignOperations) TestInsufficientMultipleItemsWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(
		nid0,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/1",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	n1 := nft.NewNFT(
		nid1,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/2",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)

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
	items := []SignItem{
		t.newSignItem(CreatorQualification, nid0, t.cid),
		t.newSignItem(CreatorQualification, nid1, t.cid),
	}
	fact := NewSignFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	signOp, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignOperations) TestInSufficientBalanceWithFee() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)

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
	items := []SignItem{
		t.newSignItem(CreatorQualification, nid, t.cid),
	}
	fact := NewSignFact(token, sender.Address, items)
	sig, err := base.NewFactSignature(sender.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig)}
	signOp, err := NewSign(fact, fs, "")
	t.NoError(err)

	err = opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignOperations) TestSameSenders() {
	sts := []state.State{}

	senderBalance := currency.NewAmount(currency.NewBig(33), t.cid)
	sender, sst := t.newAccount(true, []currency.Amount{senderBalance})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	sts = append(sts, sst...)
	sts = append(sts, pst)

	nid0 := nft.NewNFTID(t.symbol, 1)
	nid1 := nft.NewNFTID(t.symbol, 2)
	n0 := nft.NewNFT(
		nid0,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/1",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	n1 := nft.NewNFT(
		nid1,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/1",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)

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
	items0 := []SignItem{
		t.newSignItem(CreatorQualification, nid0, t.cid),
	}
	fact0 := NewSignFact(token0, sender.Address, items0)
	sig0, err := base.NewFactSignature(sender.Privs()[0], fact0, nil)
	t.NoError(err)
	fs0 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig0)}
	approve0, err := NewSign(fact0, fs0, "")
	t.NoError(err)

	t.NoError(opr.Process(approve0))

	token1 := util.UUID().Bytes()
	items1 := []SignItem{
		t.newSignItem(CreatorQualification, nid1, t.cid),
	}
	fact1 := NewSignFact(token1, sender.Address, items1)
	sig1, err := base.NewFactSignature(sender.Privs()[0], fact1, nil)
	t.NoError(err)
	fs1 := []base.FactSign{base.NewBaseFactSign(sender.Privs()[0].Publickey(), sig1)}
	approve1, err := NewSign(fact1, fs1, "")
	t.NoError(err)

	err = opr.Process(approve1)

	t.Contains(err.Error(), "violates only one sender")
}

func (t *testSignOperations) TestUnderThreshold() {
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

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender,
		"",
		"https://localhost:5000/nft/1",
		sender,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)

	nst := t.newStateNFT(n)
	_, dst := t.newCollectionDesign(true, parent, sender, []base.Address{sender}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		pst,
		nst,
	)
	sts = append(sts, dst...)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}
	signOp := t.newSign(sender, pks, items)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testSignOperations) TestUnknownKey() {
	sender, sst := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(1), t.cid)})
	parent, _, pst := t.newContractAccount(true, true, sender.Address)

	nid := nft.NewNFTID(t.symbol, 1)
	n := nft.NewNFT(
		nid,
		true,
		sender.Address,
		"",
		"https://localhost:5000/nft/1",
		sender.Address,
		nft.NewSigners(0, []nft.Signer{nft.NewSigner(sender.Address, 0, false)}),
		nft.NewSigners(0, []nft.Signer{}),
	)
	nst := t.newStateNFT(n)

	_, dst := t.newCollectionDesign(true, parent, sender.Address, []base.Address{sender.Address}, t.symbol, []nft.NFTID{nid}, []nft.NFTID{})

	sts := []state.State{}
	sts = append(sts, sst...)
	sts = append(sts, pst)
	sts = append(sts, dst...)
	sts = append(sts, nst)

	pool, _ := t.statepool(sts)
	feeer := extensioncurrency.NewFixedFeeer(sender.Address, currency.ZeroBig, currency.ZeroBig)

	cp := extensioncurrency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), nft.NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignItem{t.newSignItem(CreatorQualification, nid, t.cid)}

	signOp := t.newSign(sender.Address, []key.Privatekey{sender.Priv, key.NewBasePrivatekey()}, items)

	err := opr.Process(signOp)

	var oper operation.ReasonError
	t.True(errors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestSignOperations(t *testing.T) {
	suite.Run(t, new(testSignOperations))
}
