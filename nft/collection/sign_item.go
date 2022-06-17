package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
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
		return isvalid.InvalidError.Errorf("invalid qualification; %q", q)
	}
	return nil
}

var (
	SignItemType   = hint.Type("mitum-nft-sign-item")
	SignItemHint   = hint.NewHint(SignItemType, "v0.0.1")
	SignItemHinter = SignItem{BaseHinter: hint.NewBaseHinter(SignItemHint)}
)

type SignItem struct {
	hint.BaseHinter
	qualification Qualification
	nft           nft.NFTID
	cid           currency.CurrencyID
}

func NewSignItem(q Qualification, n nft.NFTID, cid currency.CurrencyID) SignItem {
	return SignItem{
		BaseHinter:    hint.NewBaseHinter(SignItemHint),
		qualification: q,
		nft:           n,
		cid:           cid,
	}
}

func (it SignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.qualification.Bytes(),
		it.nft.Bytes(),
		it.cid.Bytes(),
	)
}

func (it SignItem) IsValid([]byte) error {
	return isvalid.Check(nil, false, it.BaseHinter, it.qualification, it.nft)
}

func (it SignItem) Qualification() Qualification {
	return it.qualification
}

func (it SignItem) NFT() nft.NFTID {
	return it.nft
}

func (it SignItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it SignItem) Rebuild() SignItem {
	return it
}
