package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/isaac"
	"github.com/spikeekips/mitum/storage"
	"github.com/stretchr/testify/suite"
)

func MustAddress(s string) currency.Address {
	a := currency.NewAddress(s)
	if err := a.IsValid(nil); err != nil {
		panic(err)
	}
	return a
}

type account struct { // nolint: unused
	Address base.Address
	Priv    key.Privatekey
	Key     currency.BaseAccountKey
}

func (ac *account) Privs() []key.Privatekey {
	return []key.Privatekey{ac.Priv}
}

func (ac *account) Keys() currency.AccountKeys {
	keys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{ac.Key}, 100)

	return keys
}

func generateAccount() *account { // nolint: unused
	priv := key.NewBasePrivatekey()

	key, err := currency.NewBaseAccountKey(priv.Publickey(), 100)
	if err != nil {
		panic(err)
	}

	keys, err := currency.NewBaseAccountKeys([]currency.AccountKey{key}, 100)
	if err != nil {
		panic(err)
	}

	address, _ := currency.NewAddressFromKeys(keys)

	return &account{Address: address, Priv: priv, Key: key}
}

type baseTest struct { // nolint: unused
	suite.Suite
	isaac.StorageSupportTest
	cid currency.CurrencyID
}

func (t *baseTest) SetupSuite() {
	t.StorageSupportTest.SetupSuite()

	_ = t.Encs.TestAddHinter(key.BasePublickey{})
	_ = t.Encs.TestAddHinter(base.BaseFactSign{})
	_ = t.Encs.TestAddHinter(currency.AccountKeyHinter)
	_ = t.Encs.TestAddHinter(currency.AccountKeysHinter)
	_ = t.Encs.TestAddHinter(currency.AddressHinter)
	_ = t.Encs.TestAddHinter(currency.CreateAccountsHinter)
	_ = t.Encs.TestAddHinter(currency.AccountHinter)
	_ = t.Encs.TestAddHinter(CollectionRegisterHinter)
	_ = t.Encs.TestAddHinter(CollectionPolicyUpdaterHinter)
	_ = t.Encs.TestAddHinter(MintHinter)
	_ = t.Encs.TestAddHinter(TransferHinter)
	_ = t.Encs.TestAddHinter(ApproveHinter)
	_ = t.Encs.TestAddHinter(DelegateHinter)
	_ = t.Encs.TestAddHinter(BurnHinter)
	_ = t.Encs.TestAddHinter(SignHinter)

	t.cid = currency.CurrencyID("SEEME")
}

func (t *baseTest) newAccount() *account {
	return generateAccount()
}

type baseTestOperationProcessor struct { // nolint: unused
	baseTest
}

func (t *baseTestOperationProcessor) statepool(s ...[]state.State) (*storage.Statepool, prprocessor.OperationProcessor) {
	base := map[string]state.State{}
	for _, l := range s {
		for _, st := range l {
			base[st.Key()] = st
		}
	}

	pool, err := storage.NewStatepoolWithBase(t.Database(nil, nil), base)
	t.NoError(err)

	opr := (NewOperationProcessor(nil)).New(pool)

	return pool, opr
}

func (t *baseTestOperationProcessor) newStateKeys(a base.Address, keys currency.AccountKeys) state.State {
	key := currency.StateKeyAccount(a)

	ac, err := currency.NewAccount(a, keys)
	t.NoError(err)

	value, _ := state.NewHintedValue(ac)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newKey(pub key.Publickey, w uint) currency.BaseAccountKey {
	k, err := currency.NewBaseAccountKey(pub, w)
	if err != nil {
		panic(err)
	}

	return k
}

func (t *baseTestOperationProcessor) newAccount(exists bool, amounts []currency.Amount) (*account, []state.State) {
	ac := t.baseTest.newAccount()

	if !exists {
		return ac, nil
	}

	var sts []state.State
	sts = append(sts, t.newStateKeys(ac.Address, ac.Keys()))

	for _, am := range amounts {
		sts = append(sts, t.newStateAmount(ac.Address, am))
	}

	return ac, sts
}

func (t *baseTestOperationProcessor) newContractAccount(exists, active bool, owner base.Address) (base.Address, extensioncurrency.ContractAccount, state.State) {
	ks := extensioncurrency.NewContractAccountKeys()
	ac, err := currency.NewAddressFromKeys(ks)
	t.NoError(err)

	ca := extensioncurrency.NewContractAccount(owner, active)

	if !exists {
		return ac, ca, nil
	}

	key := extensioncurrency.StateKeyContractAccount(ac)
	value, _ := state.NewHintedValue(ca)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return ac, ca, su
}

func (t *baseTestOperationProcessor) newCollectionDesign(active bool, parent, creator base.Address, whites []base.Address, symbol extensioncurrency.ContractID, actives, deactives []nft.NFTID) (nft.Design, []state.State) {
	policy := NewCollectionPolicy("Collection", 0, "", whites)
	design := nft.NewDesign(parent, creator, symbol, active, policy)
	t.NoError(design.IsValid(nil))

	collectionKey := StateKeyCollection(symbol)
	collectionValue, _ := state.NewHintedValue(design)
	collectionState, err := state.NewStateV0(collectionKey, collectionValue, base.NilHeight)
	t.NoError(err)

	nftsKey := StateKeyNFTs(symbol)
	nftsBox := NewNFTBox(actives)
	nftsValue, _ := state.NewHintedValue(nftsBox)
	nftsState, err := state.NewStateV0(nftsKey, nftsValue, base.NilHeight)

	idxKey := StateKeyCollectionLastIDX(symbol)
	idxValue, _ := state.NewNumberValue(uint64(len(actives)) + uint64(len(deactives)))
	idxState, err := state.NewStateV0(idxKey, idxValue, base.NilHeight)

	sts := []state.State{collectionState, nftsState, idxState}

	return design, sts
}

func (t *baseTestOperationProcessor) newStateNFT(n nft.NFT) state.State {
	nftKey := StateKeyNFT(n.ID())
	nftValue, _ := state.NewHintedValue(n)
	st, err := state.NewStateV0(nftKey, nftValue, base.NilHeight)
	t.NoError(err)

	return st
}

func (t *baseTestOperationProcessor) newStateAmount(a base.Address, amount currency.Amount) state.State {
	key := currency.StateKeyBalance(a, amount.Currency())
	value, _ := state.NewHintedValue(amount)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateBalance(a base.Address, big currency.Big, cid currency.CurrencyID) state.State {
	key := currency.StateKeyBalance(a, cid)
	value, _ := state.NewHintedValue(currency.NewAmount(big, cid))
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateAgent(a base.Address, symbol extensioncurrency.ContractID, agents []base.Address) state.State {
	key := StateKeyAgents(a, symbol)
	box := NewAgentBox(symbol, agents)
	value, _ := state.NewHintedValue(box)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newCurrencyDesignState(cid currency.CurrencyID, big currency.Big, genesisAccount base.Address, feeer extensioncurrency.Feeer) state.State {
	de := extensioncurrency.NewCurrencyDesign(currency.NewAmount(big, cid), genesisAccount, extensioncurrency.NewCurrencyPolicy(currency.ZeroBig, feeer))

	st, err := state.NewStateV0(extensioncurrency.StateKeyCurrencyDesign(cid), nil, base.NilHeight)
	t.NoError(err)

	nst, err := extensioncurrency.SetStateCurrencyDesignValue(st, de)
	t.NoError(err)

	return nst
}
