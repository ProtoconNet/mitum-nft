package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/localtime"
)

type baseTestEncode struct {
	suite.Suite

	enc       encoder.Encoder
	encs      *encoder.Encoders
	newObject func() interface{}
	encode    func(encoder.Encoder, interface{}) ([]byte, error)
	decode    func(encoder.Encoder, []byte) (interface{}, error)
	compare   func(interface{}, interface{})
}

func (t *baseTestEncode) SetupSuite() {
	t.encs = encoder.NewEncoders()
	t.encs.AddEncoder(t.enc)

	t.encs.TestAddHinter(key.BasePublickey{})
	t.encs.TestAddHinter(base.StringAddressHinter)
	t.encs.TestAddHinter(currency.AddressHinter)
	t.encs.TestAddHinter(base.BaseFactSignHinter)
	t.encs.TestAddHinter(currency.AccountKeyHinter)
	t.encs.TestAddHinter(currency.AccountKeysHinter)
	t.encs.TestAddHinter(CollectionRegisterFactHinter)
	t.encs.TestAddHinter(CollectionRegisterFormHinter)
	t.encs.TestAddHinter(CollectionRegisterHinter)
	t.encs.TestAddHinter(CollectionPolicyUpdaterFactHinter)
	t.encs.TestAddHinter(CollectionPolicyUpdaterHinter)
	t.encs.TestAddHinter(MintFactHinter)
	t.encs.TestAddHinter(MintFormHinter)
	t.encs.TestAddHinter(MintItemHinter)
	t.encs.TestAddHinter(MintHinter)
	t.encs.TestAddHinter(currency.AccountHinter)
	t.encs.TestAddHinter(TransferFactHinter)
	t.encs.TestAddHinter(TransferItemHinter)
	t.encs.TestAddHinter(TransferHinter)
	t.encs.TestAddHinter(BurnFactHinter)
	t.encs.TestAddHinter(BurnItemHinter)
	t.encs.TestAddHinter(BurnHinter)
	t.encs.TestAddHinter(SignFactHinter)
	t.encs.TestAddHinter(nft.SignerHinter)
	t.encs.TestAddHinter(nft.SignersHinter)
	t.encs.TestAddHinter(SignItemHinter)
	t.encs.TestAddHinter(SignHinter)
	t.encs.TestAddHinter(ApproveFactHinter)
	t.encs.TestAddHinter(ApproveItemHinter)
	t.encs.TestAddHinter(ApproveHinter)
	t.encs.TestAddHinter(DelegateFactHinter)
	t.encs.TestAddHinter(DelegateItemHinter)
	t.encs.TestAddHinter(DelegateHinter)
	t.encs.TestAddHinter(nft.NFTHinter)
	t.encs.TestAddHinter(nft.NFTIDHinter)
	t.encs.TestAddHinter(nft.DesignHinter)
	t.encs.TestAddHinter(CollectionPolicyHinter)
}

func (t *baseTestEncode) TestEncode() {
	i := t.newObject()

	var err error

	var b []byte
	if t.encode != nil {
		b, err = t.encode(t.enc, i)
		t.NoError(err)
	} else {
		b, err = t.enc.Marshal(i)
		t.NoError(err)
	}

	var v interface{}
	if t.decode != nil {
		v, err = t.decode(t.enc, b)
		t.NoError(err)
	} else {
		v, err = t.enc.Decode(b)
		t.NoError(err)
	}

	t.compare(i, v)
}

type baseTestOperationEncode struct {
	baseTestEncode
}

func (t *baseTestOperationEncode) TestEncode() {
	i := t.newObject()
	op, ok := i.(operation.Operation)
	t.True(ok)

	b, err := t.enc.Marshal(op)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)

	uop, ok := hinter.(operation.Operation)
	t.True(ok)

	fact := op.Fact().(operation.OperationFact)
	ufact := uop.Fact().(operation.OperationFact)
	t.True(fact.Hash().Equal(ufact.Hash()))
	t.True(fact.Hint().Equal(ufact.Hint()))
	t.Equal(fact.Token(), ufact.Token())

	t.True(op.Hash().Equal(uop.Hash()))

	t.Equal(len(op.Signs()), len(uop.Signs()))
	for i := range op.Signs() {
		a := op.Signs()[i]
		b := uop.Signs()[i]

		t.True(a.Signer().Equal(b.Signer()))
		t.Equal(a.Signature(), b.Signature())
		t.True(localtime.Equal(a.SignedAt(), b.SignedAt()))
	}

	t.compare(op, uop)
}
