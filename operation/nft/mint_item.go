package nft

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CollectionItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var MintItemHint = hint.MustNewHint("mitum-nft-mint-item-v0.0.1")

type MintItem struct {
	hint.BaseHinter
	contract mitumbase.Address
	receiver mitumbase.Address
	hash     types.NFTHash
	uri      types.URI
	creators types.Signers
	currency currencytypes.CurrencyID
}

func NewMintItem(
	contract mitumbase.Address,
	receiver mitumbase.Address,
	hash types.NFTHash,
	uri types.URI,
	creators types.Signers,
	currency currencytypes.CurrencyID,
) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		contract:   contract,
		receiver:   receiver,
		hash:       hash,
		uri:        uri,
		creators:   creators,
		currency:   currency,
	}
}

func (it MintItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.hash.Bytes(),
		it.uri.Bytes(),
		it.creators.Bytes(),
		it.currency.Bytes(),
	)
}

func (it MintItem) IsValid([]byte) error {
	return util.CheckIsValiders(
		nil,
		false,
		it.BaseHinter,
		it.contract,
		it.receiver,
		it.hash,
		it.uri,
		it.creators,
		it.currency,
	)
}

func (it MintItem) Contract() mitumbase.Address {
	return it.contract
}

func (it MintItem) Receiver() mitumbase.Address {
	return it.receiver
}

func (it MintItem) NFTHash() types.NFTHash {
	return it.hash
}

func (it MintItem) URI() types.URI {
	return it.uri
}

func (it MintItem) Creators() types.Signers {
	return it.creators
}

func (it MintItem) Addresses() ([]mitumbase.Address, error) {
	as := []mitumbase.Address{}
	as = append(as, it.receiver)
	as = append(as, it.creators.Addresses()...)

	return as, nil
}

func (it MintItem) Currency() currencytypes.CurrencyID {
	return it.currency
}
