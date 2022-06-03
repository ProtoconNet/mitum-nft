package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form *CollectionRegisterForm) unpack(
	enc encoder.Encoder,
	bTarget base.AddressDecoder,
	symbol string,
	name string,
	royalty uint,
	uri string,
) error {
	target, err := bTarget.Encode(enc)
	if err != nil {
		return err
	}
	form.target = target

	form.symbol = extensioncurrency.ContractID(symbol)
	form.name = CollectionName(name)
	form.royalty = nft.PaymentParameter(royalty)
	form.uri = nft.URI(uri)

	return nil
}

func (fact *CollectionRegisterFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bForm []byte,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	var form CollectionRegisterForm
	if hinter, err := enc.Decode(bForm); err != nil {
		return err
	} else if i, ok := hinter.(CollectionRegisterForm); !ok {
		return util.WrongTypeError.Errorf("not CollectionRegisterForm; %T", hinter)
	} else {
		form = i
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.form = form
	fact.cid = currency.CurrencyID(cid)

	return nil
}
