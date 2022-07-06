package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MintFormType   = hint.Type("mitum-nft-mint-form")
	MintFormHint   = hint.NewHint(MintFormType, "v0.0.1")
	MintFormHinter = MintForm{BaseHinter: hint.NewBaseHinter(MintFormHint)}
)

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

func MustNewMintform(hash nft.NFTHash, uri nft.URI, creators nft.Signers, copyrighters nft.Signers) MintForm {
	form := NewMintForm(hash, uri, creators, copyrighters)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form MintForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.hash.Bytes(),
		[]byte(form.uri.String()),
		form.creators.Bytes(),
		form.copyrighters.Bytes(),
	)
}

func (form MintForm) NftHash() nft.NFTHash {
	return form.hash
}

func (form MintForm) Uri() nft.URI {
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

	if ads, err := form.creators.Addresses(); err != nil {
		as = append(as, ads...)
	}

	if ads, err := form.copyrighters.Addresses(); err != nil {
		as = append(as, ads...)
	}

	return as, nil
}

func (form MintForm) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.hash,
		form.uri,
		form.creators,
		form.copyrighters); err != nil {
		return err
	}

	if len(form.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	return nil
}

var (
	MintItemType   = hint.Type("mitum-nft-mint-item")
	MintItemHint   = hint.NewHint(MintItemType, "v0.0.1")
	MintItemHinter = MintItem{BaseHinter: hint.NewBaseHinter(MintItemHint)}
)

type MintItem struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	form       MintForm
	cid        currency.CurrencyID
}

func NewMintItem(symbol extensioncurrency.ContractID, form MintForm, cid currency.CurrencyID) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		collection: symbol,
		form:       form,
		cid:        cid,
	}
}

func (it MintItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.collection.Bytes(),
		it.form.Bytes(),
		it.cid.Bytes(),
	)
}

func (it MintItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.collection, it.form, it.cid); err != nil {
		return err
	}

	return nil
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
	return it.cid
}

func (it MintItem) Rebuild() MintItem {
	return it
}
