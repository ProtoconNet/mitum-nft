package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

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
	hash        nft.NFTHash
	uri         nft.URI
	copyrighter base.Address
}

func NewMintForm(hash nft.NFTHash, uri nft.URI, copyrighter base.Address) MintForm {
	return MintForm{
		BaseHinter:  hint.NewBaseHinter(MintFormHint),
		hash:        hash,
		uri:         uri,
		copyrighter: copyrighter,
	}
}

func MustNewMintform(hash nft.NFTHash, uri nft.URI, copyrighter base.Address) MintForm {
	form := NewMintForm(hash, uri, copyrighter)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form MintForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.hash.Bytes(),
		[]byte(form.uri.String()),
		form.copyrighter.Bytes(),
	)
}

func (form MintForm) NftHash() nft.NFTHash {
	return form.hash
}

func (form MintForm) Uri() nft.URI {
	return form.uri
}

func (form MintForm) Copyrighter() base.Address {
	return form.copyrighter
}

func (form MintForm) IsValid([]byte) error {
	if err := form.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if len(form.uri.String()) < 1 {
		return isvalid.InvalidError.Errorf("empty uri")
	}

	if len(form.copyrighter.String()) > 0 {
		if err := form.copyrighter.IsValid(nil); err != nil {
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

	foundHash := map[nft.NFTHash]bool{}
	for i := range it.forms {
		if err := it.forms[i].IsValid(nil); err != nil {
			return err
		}
		h := it.forms[i].NftHash()
		if _, found := foundHash[h]; found {
			return errors.Errorf("duplicated nft hash found; %s", h)
		}
		foundHash[h] = true
	}

	return nil
}

func (it BaseMintItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it BaseMintItem) Addresses() []base.Address {
	as := []base.Address{}

	for i := range it.forms {
		if len(it.forms[i].Copyrighter().String()) > 0 {
			as = append(as, it.forms[i].Copyrighter())
		}
	}

	return as
}

func (it BaseMintItem) Hashes() []nft.NFTHash {
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
