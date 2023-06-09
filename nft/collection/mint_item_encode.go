package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *MintItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col, hs, uri string,
	bcr []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal MintItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.collection = currencybase.ContractID(col)
	it.hash = nft.NFTHash(hs)
	it.uri = nft.URI(uri)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	if hinter, err := enc.Decode(bcr); err != nil {
		return e(err, "")
	} else if creators, ok := hinter.(nft.Signers); !ok {
		return e(util.ErrWrongType.Errorf("expected Signers, not %T", hinter), "")
	} else {
		it.creators = creators
	}

	it.currency = currencybase.CurrencyID(cid)

	return nil
}
