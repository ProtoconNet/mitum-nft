package collection

import (
	"strings"
	"testing"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type testTransfer struct {
	suite.Suite
}

func (t *testTransfer) TestNew() {
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

func (t *testTransfer) TestDeplicateNFTID() {
	sender := MustAddress(util.UUID().String())
	receiver := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	nid := nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1)

	items := []TransferItem{
		NewTransferItem(receiver, nid, "MCC"),
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
	t.Contains(err.Error(), "duplicate nft found")
}

func (t *testTransfer) TestEmptyItems() {
	sender := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []TransferItem{}
	fact := NewTransferFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	transfer, err := NewTransfer(fact, fs, "")
	t.NoError(err)

	err = transfer.IsValid(nil)
	t.Contains(err.Error(), "empty items for TransferFact")
}

func (t *testTransfer) TestOverMaxItems() {
	sender := MustAddress(util.UUID().String())
	receiver := MustAddress(util.UUID().String())
	token := util.UUID().Bytes()

	items := []TransferItem{
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 1), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 2), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 3), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 4), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 5), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 6), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 7), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 8), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 9), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 10), "MCC"),
		NewTransferItem(receiver, nft.NewNFTID(extensioncurrency.ContractID("ABC"), 11), "MCC"),
	}
	fact := NewTransferFact(token, sender, items)

	pk := key.NewBasePrivatekey()
	sig, err := base.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []base.FactSign{base.NewBaseFactSign(pk.Publickey(), sig)}

	transfer, err := NewTransfer(fact, fs, "")
	t.NoError(err)

	err = transfer.IsValid(nil)
	t.Contains(err.Error(), "items over allowed")
}

func (t *testTransfer) TestOverSizeMemo() {
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

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	transfer, err := NewTransfer(fact, fs, memo)
	t.NoError(err)

	err = transfer.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestTransfers(t *testing.T) {
	suite.Run(t, new(testTransfer))
}
