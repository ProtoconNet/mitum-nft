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
	creators     []nft.RightHolder
	copyrighters []nft.RightHolder
}

func NewMintForm(hash nft.NFTHash, uri nft.URI, creators []nft.RightHolder, copyrighters []nft.RightHolder) MintForm {
	return MintForm{
		BaseHinter:   hint.NewBaseHinter(MintFormHint),
		hash:         hash,
		uri:          uri,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func MustNewMintform(hash nft.NFTHash, uri nft.URI, creators []nft.RightHolder, copyrighters []nft.RightHolder) MintForm {
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

func (form MintForm) Creators() []nft.RightHolder {
	return form.creators
}

func (form MintForm) Copyrighters() []nft.RightHolder {
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
	if err := form.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if len(form.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	for i := range form.creators {
		if err := form.creators[i].IsValid(nil); err != nil {
			return err
		}
	}

	for i := range form.copyrighters {
		if err := form.copyrighters[i].IsValid(nil); err != nil {
			return err
		}
	}

	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.hash); err != nil {
		return err
	}

	return nil
}

type BaseMintItem struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	forms      []MintForm
	cid        currency.CurrencyID
}

func NewBaseMintItem(ht hint.Hint, collection extensioncurrency.ContractID, forms []MintForm, cid currency.CurrencyID) BaseMintItem {
	return BaseMintItem{
		BaseHinter: hint.NewBaseHinter(ht),
		collection: collection,
		forms:      forms,
		cid:        cid,
	}
}

func (it BaseMintItem) Bytes() []byte {
	bf := make([][]byte, len(it.forms))

	for i := range it.forms {
		bf[i] = it.forms[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.collection.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(bf...),
	)
}

func (it BaseMintItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.collection, it.cid); err != nil {
		return err
	}

	if len(it.forms) < 1 {
		return isvalid.InvalidError.Errorf("empty forms for BaseMintItem")
	}

	for i := range it.forms {
		if err := it.forms[i].IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (it BaseMintItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it BaseMintItem) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	for i := range it.forms {
		if adr, err := it.forms[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, adr...)
		}
	}

	return as, nil
}

func (it BaseMintItem) NftHashes() []nft.NFTHash {
	hs := make([]nft.NFTHash, len(it.forms))

	for i := range it.forms {
		hs[i] = it.forms[i].NftHash()
	}

	return hs
}

func (it BaseMintItem) Forms() []MintForm {
	return it.forms
}

func (it BaseMintItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseMintItem) Rebuild() MintItem {
	forms := make([]MintForm, len(it.forms))
	for i := range it.forms {
		forms[i] = it.forms[i]
	}
	it.forms = forms

	return it
}
