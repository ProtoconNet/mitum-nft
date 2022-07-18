package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *CollectionPolicyUpdaterFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bs base.AddressDecoder,
	collection string,
	bp []byte,
	cid string,
) error {
	sender, err := bs.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(bp); err != nil {
		return err
	} else if p, ok := hinter.(CollectionPolicy); !ok {
		return util.WrongTypeError.Errorf("not CollectionPolicy; %T", hinter)
	} else {
		fact.policy = p
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.collection = extensioncurrency.ContractID(collection)
	fact.cid = currency.CurrencyID(cid)

	return nil
}
