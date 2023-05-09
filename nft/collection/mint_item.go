package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MintFormHint = hint.MustNewHint("mitum-nft-mint-form-v0.0.1")

type MintForm struct {
	hint.BaseHinter
	hash         nft.NFTHash
	uri          nft.URI
	creators     nft.Signers
	copyrighters nft.Signers
}

func NewMintForm(hash nft.NFTHash, uri nft.URI, creators nft.Signers, copyrighters nft.Signers) MintForm {
	return MintForm{
		BaseHinter:   hint.NewBaseHinter(MintFormHint),
		hash:         hash,
		uri:          uri,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func (form MintForm) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		form.BaseHinter,
		form.hash,
		form.uri,
		form.creators,
		form.copyrighters,
	); err != nil {
		return err
	}

	if len(form.uri.String()) < 1 {
		return util.ErrInvalid.Errorf("empty uri")
	}

	return nil
}

func (form MintForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.hash.Bytes(),
		form.uri.Bytes(),
		form.creators.Bytes(),
		form.copyrighters.Bytes(),
	)
}

func (form MintForm) NFTHash() nft.NFTHash {
	return form.hash
}

func (form MintForm) URI() nft.URI {
	return form.uri
}

func (form MintForm) Creators() nft.Signers {
	return form.creators
}

func (form MintForm) Copyrighters() nft.Signers {
	return form.copyrighters
}

func (form MintForm) Addresses() ([]base.Address, error) {
	as := []base.Address{}
	as = append(as, form.creators.Addresses()...)
	as = append(as, form.copyrighters.Addresses()...)

	return as, nil
}

type CollectionItem interface {
	util.Byter
	util.IsValider
	Currency() currency.CurrencyID
}

var MintItemHint = hint.MustNewHint("mitum-nft-mint-item-v0.0.1")

type MintItem struct {
	hint.BaseHinter
	contract   base.Address
	collection extensioncurrency.ContractID
	form       MintForm
	currency   currency.CurrencyID
}

func NewMintItem(contract base.Address, collection extensioncurrency.ContractID, form MintForm, currency currency.CurrencyID) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		contract:   contract,
		collection: collection,
		form:       form,
		currency:   currency,
	}
}

func (it MintItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.form.Bytes(),
		it.currency.Bytes(),
	)
}

func (it MintItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.collection, it.form, it.currency)
}

func (it MintItem) Contract() base.Address {
	return it.contract
}

func (it MintItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it MintItem) Addresses() ([]base.Address, error) {
	return it.form.Addresses()
}

func (it MintItem) Form() MintForm {
	return it.form
}

func (it MintItem) Currency() currency.CurrencyID {
	return it.currency
}
