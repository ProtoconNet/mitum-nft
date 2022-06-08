package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (form *MintForm) unpack(
	enc encoder.Encoder,
	hash string,
	uri string,
	_copyrighter string,
) error {
	form.hash = nft.NFTHash(hash)
	form.uri = nft.URI(uri)

	if len(_copyrighter) < 1 {
		form.copyrighter = currency.Address{}
	} else {
		copyrighter, err := base.DecodeAddressFromString(_copyrighter, enc)
		if err != nil {
			return err
		}
		form.copyrighter = copyrighter
	}

	return nil
}

func (it *BaseMintItem) unpack(
	enc encoder.Encoder,
	collection string,
	bForms []byte,
	cid string,
) error {
	it.collection = extensioncurrency.ContractID(collection)

	hForms, err := enc.DecodeSlice(bForms)
	if err != nil {
		return err
	}

	forms := make([]MintForm, len(hForms))
	for i := range hForms {
		j, ok := hForms[i].(MintForm)
		if !ok {
			return util.WrongTypeError.Errorf("not MintForm; %T", hForms[i])
		}
		forms[i] = j
	}

	it.forms = forms
	it.cid = currency.CurrencyID(cid)

	return nil
}
