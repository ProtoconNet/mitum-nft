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
	bt base.AddressDecoder,
	symbol string,
	name string,
	royalty uint,
	uri string,
) error {
	target, err := bt.Encode(enc)
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
	bs base.AddressDecoder,
	bf []byte,
	cid string,
) error {
	sender, err := bs.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(bf); err != nil {
		return err
	} else if form, ok := hinter.(CollectionRegisterForm); !ok {
		return util.WrongTypeError.Errorf("not CollectionRegisterForm; %T", hinter)
	} else {
		fact.form = form
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.cid = currency.CurrencyID(cid)

	return nil
}
