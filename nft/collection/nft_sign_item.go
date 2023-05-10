package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	CreatorQualification     = Qualification("creator")
	CopyrighterQualification = Qualification("copyrighter")
)

type Qualification string

func (q Qualification) Bytes() []byte {
	return []byte(q)
}

func (q Qualification) String() string {
	return string(q)
}

func (q Qualification) IsValid([]byte) error {
	if !(q == CreatorQualification || q == CopyrighterQualification) {
		return util.ErrInvalid.Errorf("invalid qualification, %q", q)
	}
	return nil
}

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	contract      base.Address
	collection    extensioncurrency.ContractID
	qualification Qualification
	nft           nft.NFTID
	currency      currency.CurrencyID
}

func NewNFTSignItem(contract base.Address, collection extensioncurrency.ContractID, q Qualification, n nft.NFTID, currency currency.CurrencyID) NFTSignItem {
	return NFTSignItem{
		BaseHinter:    hint.NewBaseHinter(NFTSignItemHint),
		contract:      contract,
		collection:    collection,
		qualification: q,
		nft:           n,
		currency:      currency,
	}
}

func (it NFTSignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.qualification.Bytes(),
		it.nft.Bytes(),
		it.currency.Bytes(),
	)
}

func (it NFTSignItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.collection, it.qualification, it.nft, it.currency)
}

func (it NFTSignItem) Qualification() Qualification {
	return it.qualification
}

func (it NFTSignItem) NFT() nft.NFTID {
	return it.nft
}

func (it NFTSignItem) Contract() base.Address {
	return it.contract
}

func (it NFTSignItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it NFTSignItem) Collection() extensioncurrency.ContractID {
	return it.collection
}
