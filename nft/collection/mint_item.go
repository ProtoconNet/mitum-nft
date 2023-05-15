package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

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
	hash       nft.NFTHash
	uri        nft.URI
	creators   nft.Signers
	currency   currency.CurrencyID
}

func NewMintItem(
	contract base.Address,
	collection extensioncurrency.ContractID,
	hash nft.NFTHash,
	uri nft.URI,
	creators nft.Signers,
	currency currency.CurrencyID,
) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		contract:   contract,
		collection: collection,
		hash:       hash,
		uri:        uri,
		creators:   creators,
		currency:   currency,
	}
}

func (it MintItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.hash.Bytes(),
		it.uri.Bytes(),
		it.creators.Bytes(),
		it.currency.Bytes(),
	)
}

func (it MintItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.collection, it.hash, it.uri, it.creators, it.currency)
}

func (it MintItem) Contract() base.Address {
	return it.contract
}

func (it MintItem) Collection() extensioncurrency.ContractID {
	return it.collection
}

func (it MintItem) NFTHash() nft.NFTHash {
	return it.hash
}

func (it MintItem) URI() nft.URI {
	return it.uri
}

func (it MintItem) Creators() nft.Signers {
	return it.creators
}

func (it MintItem) Addresses() ([]base.Address, error) {
	as := []base.Address{}
	as = append(as, it.creators.Addresses()...)

	return as, nil
}

func (it MintItem) Currency() currency.CurrencyID {
	return it.currency
}
