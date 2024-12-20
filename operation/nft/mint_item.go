package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

type CollectionItem interface {
	util.Byter
	util.IsValider
	Currency() ctypes.CurrencyID
}

var MintItemHint = hint.MustNewHint("mitum-nft-mint-item-v0.0.1")

type MintItem struct {
	hint.BaseHinter
	contract base.Address
	receiver base.Address
	hash     types.NFTHash
	uri      types.URI
	creators types.Signers
	currency ctypes.CurrencyID
}

func NewMintItem(
	contract base.Address,
	receiver base.Address,
	hash types.NFTHash,
	uri types.URI,
	creators types.Signers,
	currency ctypes.CurrencyID,
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
	if it.receiver.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", it.receiver))
	}

	signers := it.creators.Signers()
	for _, signer := range signers {
		if signer.Address().Equal(it.contract) {
			return common.ErrSelfTarget.Wrap(errors.Errorf("creator %v is same with contract account", signer.Address()))
		}

		if signer.Signed() {
			return common.ErrValueInvalid.Wrap(errors.Errorf("creator %v should not be signed at the time of minting", signer.Address()))
		}
	}

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

func (it MintItem) Contract() base.Address {
	return it.contract
}

func (it MintItem) Receiver() base.Address {
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

func (it MintItem) Addresses() ([]base.Address, error) {
	as := []base.Address{}
	as = append(as, it.receiver)
	as = append(as, it.creators.Addresses()...)

	return as, nil
}

func (it MintItem) Currency() ctypes.CurrencyID {
	return it.currency
}
