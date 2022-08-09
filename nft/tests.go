package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	TestPolicyType   = hint.Type("mitum-nft-test-policy")
	TestPolicyHint   = hint.NewHint(TestPolicyType, "v0.0.1")
	TestPolicyHinter = TestPolicy{BaseHinter: hint.NewBaseHinter(TestPolicyHint)}
)

type TestPolicy struct {
	hint.BaseHinter
	value int
}

func (tp TestPolicy) IsValid([]byte) error {
	if _, ok := interface{}(tp).(BasePolicy); !ok {
		return errors.Errorf("not BasePolicy; %T", tp)
	}
	return nil
}

func (tp TestPolicy) Bytes() []byte {
	b := make([]byte, 0)
	return b
}

func (tp TestPolicy) Addresses() ([]base.Address, error) {
	ads := make([]base.Address, 0)
	return ads, nil
}

func (tp TestPolicy) Equal(c BasePolicy) bool {
	cp, ok := c.(TestPolicy)
	if !ok {
		return false
	}

	return tp.value == cp.value
}

func (tp TestPolicy) Rebuild() BasePolicy {
	return tp
}

func (tp TestPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(tp.Hint()),
		bson.M{
			"value": tp.value,
		}),
	)
}

type TestPolicyBSONUnpacker struct {
	VL int `bson:"value"`
}

func (tp *TestPolicy) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var utp TestPolicyBSONUnpacker
	if err := enc.Unmarshal(b, &utp); err != nil {
		return err
	}

	return tp.unpack(enc, utp.VL)
}

func (tp *TestPolicy) unpack(
	enc encoder.Encoder,
	v int,
) error {
	tp.value = v
	return nil
}

type TestPolicyJSONPacker struct {
	jsonenc.HintedHead
	VL int `json:"value"`
}

func (tp TestPolicy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(TestPolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(tp.Hint()),
		VL:         tp.value,
	})
}

type TestPolicyJSONUnpacker struct {
	VL int `json:"value"`
}

func (tp *TestPolicy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var utp TestPolicyJSONUnpacker
	if err := enc.Unmarshal(b, &utp); err != nil {
		return err
	}

	return tp.unpack(enc, utp.VL)
}

func NewTestPolicy(value int) TestPolicy {
	p := TestPolicy{
		BaseHinter: hint.NewBaseHinter(TestPolicyHint),
		value:      value,
	}
	if err := p.IsValid(nil); err != nil {
		panic(err)
	}
	return p
}

func NewTestNFTID(idx uint64) NFTID {
	return MustNewNFTID(extensioncurrency.ContractID("ABC"), idx)
}

func NewTestSigners() Signers {
	signer := MustNewSigner(NewTestAddress(), 100, false)
	signers := MustNewSigners(100, []Signer{signer})
	return signers
}

func NewTestAddress() base.Address {
	k, err := currency.NewBaseAccountKey(key.NewBasePrivatekey().Publickey(), 100)
	if err != nil {
		panic(err)
	}
	if err = k.IsValid(nil); err != nil {
		panic(err)
	}

	keys, err := currency.NewBaseAccountKeys([]currency.AccountKey{k}, 100)
	if err != nil {
		panic(err)
	}
	if err = keys.IsValid(nil); err != nil {
		panic(err)
	}

	a, err := currency.NewAddressFromKeys(keys)
	if err != nil {
		panic(err)
	}
	if err = a.IsValid(nil); err != nil {
		panic(err)
	}

	return a
}
