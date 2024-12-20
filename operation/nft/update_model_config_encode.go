package nft

import (
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *UpdateModelConfigFact) unpack(
	enc encoder.Encoder,
	sd string,
	ct string,
	nm string,
	ry uint,
	uri string,
	bws []string,
	cid string,
) error {
	fact.currency = ctypes.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	contract, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.contract = contract

	fact.name = types.CollectionName(nm)
	fact.royalty = types.PaymentParameter(ry)
	fact.uri = types.URI(uri)

	switch a, err := base.DecodeAddress(ct, enc); {
	case err != nil:
		return err
	default:
		fact.contract = a
	}

	whitelist := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return err
		}
		whitelist[i] = white
	}
	fact.whitelist = whitelist

	return nil
}
