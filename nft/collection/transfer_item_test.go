package collection

import (
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testTransferItem struct {
	suite.Suite
}

func (t *testTransferItem) TestNew() {
	sender := MustAddress(util.UUID().String())
	receiver := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
	items := []TransferItem{NewTransferItem(receiver, nid, "MCC")}
	fact := NewTransferFact(token, sender, items)

	var fs []base.FactSign

	for _, pk := range []key.Privatekey{
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
		key.NewBasePrivatekey(),
	} {
		sig, err := base.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	transfer, err := NewTransfer(fact, fs, "")
	t.NoError(err)

	t.NoError(transfer.IsValid(nil))

	t.Implements((*base.Fact)(nil), transfer.Fact())
	t.Implements((*operation.Operation)(nil), transfer)
}

func (t *testTransferItem) TestZeroIDX() {
	sender := MustAddress(util.UUID().String())
	receiver := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 0)

	items := []TransferItem{
		NewTransferItem(receiver, nid, "MCC"),
	}
	fact := NewTransferFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	transfer, err := NewTransfer(fact, fs, "")
	t.NoError(err)

	err = transfer.IsValid(nil)
	t.Contains(err.Error(), "nid idx must be over zero")
}

func (t *testTransferItem) TestOverMaxIDX() {
	sender := MustAddress(util.UUID().String())
	receiver := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), uint64(nft.MaxNFTIdx+1))

	items := []TransferItem{
		NewTransferItem(receiver, nid, "MCC"),
	}
	fact := NewTransferFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	transfer, err := NewTransfer(fact, fs, "")
	t.NoError(err)

	err = transfer.IsValid(nil)
	t.Contains(err.Error(), "nid idx over max")
}

func TestTransferItem(t *testing.T) {
	suite.Run(t, new(testTransferItem))
}

func testTransferItemEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sender := MustAddress(util.UUID().String())
		receiver := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		nid0 := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)
		nid1 := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 2)
		items := []TransferItem{
			NewTransferItem(receiver, nid0, "MCC"),
			NewTransferItem(receiver, nid1, "MCC"),
		}
		fact := NewTransferFact(token, sender, items)

		var fs []base.FactSign

		for _, pk := range []key.Privatekey{
			key.NewBasePrivatekey(),
			key.NewBasePrivatekey(),
			key.NewBasePrivatekey(),
		} {
			sig, err := base.NewFactSignature(pk, fact, nil)
			t.NoError(err)

			fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
		}

		transfer, err := NewTransfer(fact, fs, "")
		t.NoError(err)

		return transfer
	}

	t.compare = func(a, b interface{}) {
		ta := a.(Transfer)
		tb := b.(Transfer)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(TransferFact)
		ufact := tb.Fact().(TransferFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.Receiver().Equal(b.Receiver()))
			t.True(a.NFT().Equal(b.NFT()))
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestTransferItemEncodeJSON(t *testing.T) {
	suite.Run(t, testTransferItemEncode(jsonenc.NewEncoder()))
}

func TestTransferItemEncodeBSON(t *testing.T) {
	suite.Run(t, testTransferItemEncode(bsonenc.NewEncoder()))
}
