package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (it *MintItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, ra, hs, uri string,
	bcr []byte,
	cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.hash = types.NFTHash(hs)
	it.uri = types.URI(uri)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(ra, enc); {
	case err != nil:
		return err
	default:
		it.receiver = a
	}

	if hinter, err := enc.Decode(bcr); err != nil {
		return err
	} else if creators, ok := hinter.(types.Signers); !ok {
		return common.ErrTypeMismatch.Wrap(errors.Errorf("expected Signers, not %T", hinter))
	} else {
		it.creators = creators
	}

	it.currency = ctypes.CurrencyID(cid)

	return nil
}
