package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (form *MintForm) unpack(
	enc encoder.Encoder,
	hash string,
	uri string,
	bcrs []byte,
	bcps []byte,
) error {
	form.hash = nft.NFTHash(hash)
	form.uri = nft.URI(uri)

	if hinter, err := enc.Decode(bcrs); err != nil {
		return err
	} else if sns, ok := hinter.(nft.Signers); !ok {
		return util.WrongTypeError.Errorf("not Signers; %T", hinter)
	} else {
		form.creators = sns
	}

	if hinter, err := enc.Decode(bcps); err != nil {
		return err
	} else if sns, ok := hinter.(nft.Signers); !ok {
		return util.WrongTypeError.Errorf("not Signer; %T", hinter)
	} else {
		form.copyrighters = sns
	}

	return nil
}

func (it *MintItem) unpack(
	enc encoder.Encoder,
	collection string,
	bf []byte,
	cid string,
) error {
	it.collection = extensioncurrency.ContractID(collection)

	if hinter, err := enc.Decode(bf); err != nil {
		return err
	} else if form, ok := hinter.(MintForm); !ok {
		return util.WrongTypeError.Errorf("not MintForm; %T", hinter)
	} else {
		it.form = form
	}

	it.cid = currency.CurrencyID(cid)

	return nil
}
