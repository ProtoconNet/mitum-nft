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

	hcrs, err := enc.DecodeSlice(bcrs)
	if err != nil {
		return err
	}
	crs := make([]nft.RightHolder, len(hcrs))
	for i := range hcrs {
		r, ok := hcrs[i].(nft.RightHolder)
		if !ok {
			return util.WrongTypeError.Errorf("not RightHolder; %T", hcrs[i])
		}
		crs[i] = r
	}
	form.creators = crs

	hcps, err := enc.DecodeSlice(bcps)
	if err != nil {
		return err
	}
	cps := make([]nft.RightHolder, len(hcps))
	for i := range hcps {
		r, ok := hcps[i].(nft.RightHolder)
		if !ok {
			return util.WrongTypeError.Errorf("not RightHolder; %T", hcps[i])
		}
		cps[i] = r
	}
	form.copyrighters = cps

	return nil
}

func (it *BaseMintItem) unpack(
	enc encoder.Encoder,
	collection string,
	bfs []byte,
	cid string,
) error {
	it.collection = extensioncurrency.ContractID(collection)

	hfs, err := enc.DecodeSlice(bfs)
	if err != nil {
		return err
	}

	forms := make([]MintForm, len(hfs))
	for i := range hfs {
		form, ok := hfs[i].(MintForm)
		if !ok {
			return util.WrongTypeError.Errorf("not MintForm; %T", hfs[i])
		}
		forms[i] = form
	}
	it.forms = forms

	it.cid = currency.CurrencyID(cid)

	return nil
}
