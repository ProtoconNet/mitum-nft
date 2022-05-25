package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *CollectionRegisterFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bDesign []byte,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	var design nft.Design
	if hinter, err := enc.Decode(bDesign); err != nil {
		return err
	} else if i, ok := hinter.(nft.Design); !ok {
		return util.WrongTypeError.Errorf("not Design; %T", hinter)
	} else {
		design = i
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.design = design
	fact.cid = currency.CurrencyID(cid)

	return nil
}
