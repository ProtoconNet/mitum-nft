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
	creators     []nft.Signer
	copyrighters []nft.Signer
}

func NewMintForm(hash nft.NFTHash, uri nft.URI, creators []nft.Signer, copyrighters []nft.Signer) MintForm {
	return MintForm{
		BaseHinter:   hint.NewBaseHinter(MintFormHint),
		hash:         hash,
		uri:          uri,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func MustNewMintform(hash nft.NFTHash, uri nft.URI, creators []nft.Signer, copyrighters []nft.Signer) MintForm {
	form := NewMintForm(hash, uri, creators, copyrighters)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form MintForm) Bytes() []byte {
	bcrs := [][]byte{}
	bcps := [][]byte{}

	for i := range form.creators {
		bcrs = append(bcrs, form.creators[i].Bytes())
	}

	for i := range form.copyrighters {
		bcps = append(bcps, form.copyrighters[i].Bytes())
	}

	return util.ConcatBytesSlice(
		form.hash.Bytes(),
		[]byte(form.uri.String()),
		util.ConcatBytesSlice(bcrs...),
		util.ConcatBytesSlice(bcps...),
	)
}

func (form MintForm) NftHash() nft.NFTHash {
	return form.hash
}

func (form MintForm) Uri() nft.URI {
	return form.uri
}

func (form MintForm) Creators() []nft.Signer {
	return form.creators
}

func (form MintForm) Copyrighters() []nft.Signer {
	return form.copyrighters
}

func (form MintForm) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	if len(form.creators) > 1 {
		for i := range form.creators {
			as = append(as, form.creators[i].Account())
		}
	}

	if len(form.copyrighters) > 1 {
		for i := range form.copyrighters {
			as = append(as, form.copyrighters[i].Account())
		}
	}

	return as, nil
}

func (form MintForm) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.hash,
		form.uri); err != nil {
		return err
	}

	if len(form.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	if l := len(form.creators); l > nft.MaxCreators {
		return isvalid.InvalidError.Errorf("creators over allowed; %d > %d", l, nft.MaxCreators)
	}

	if l := len(form.copyrighters); l > nft.MaxCopyrighters {
		return isvalid.InvalidError.Errorf("copyrighters over allowed; %d > %d", l, nft.MaxCopyrighters)
	}

	foundSigner := map[base.Address]bool{}
	for i := range form.creators {
		creator := form.creators[i].Account()
		if err := creator.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundSigner[creator]; found {
			return isvalid.InvalidError.Errorf("duplicate creator found; %q", creator)
		}

		foundSigner[creator] = true
	}

	foundSigner = map[base.Address]bool{}
	for i := range form.copyrighters {
		copyrighter := form.copyrighters[i].Account()
		if err := copyrighter.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundSigner[copyrighter]; found {
			return isvalid.InvalidError.Errorf("duplicate copyrighter found; %q", copyrighter)
		}

		foundSigner[copyrighter] = true
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

func NewMintItem(collection extensioncurrency.ContractID, form MintForm, cid currency.CurrencyID) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		collection: collection,
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
