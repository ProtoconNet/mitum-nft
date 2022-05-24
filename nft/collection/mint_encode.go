package collection

import (
	"net/url"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form *MintForm) unpack(
	enc encoder.Encoder,
	hash string,
	_uri string,
	bCopyrighter base.AddressDecoder,
) error {
	form.hash = nft.NFTHash(hash)

	if uri, err := url.Parse(_uri); err != nil {
		return err
	} else {
		form.uri = *uri
	}

	copyrighter, err := bCopyrighter.Encode(enc)
	if err != nil {
		return err
	}
	form.copyrighter = copyrighter

	return nil
}

func (fact *MintFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	collection string,
	bForm []byte,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(bForm); err != nil {
		return err
	} else if form, ok := hinter.(MintForm); !ok {
		return errors.Errorf("not MintForm; %T", hinter)
	} else {
		fact.form = form
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.collection = extension.ContractID(collection)
	fact.cid = currency.CurrencyID(cid)

	return nil
}
