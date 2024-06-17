package nft

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
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
	fact.currency = currencytypes.CurrencyID(cid)

	sender, err := mitumbase.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	contract, err := mitumbase.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.contract = contract

	fact.name = types.CollectionName(nm)
	fact.royalty = types.PaymentParameter(ry)
	fact.uri = types.URI(uri)

	switch a, err := mitumbase.DecodeAddress(ct, enc); {
	case err != nil:
		return err
	default:
		fact.contract = a
	}

	whitelist := make([]mitumbase.Address, len(bws))
	for i, bw := range bws {
		white, err := mitumbase.DecodeAddress(bw, enc)
		if err != nil {
			return err
		}
		whitelist[i] = white
	}
	fact.whitelist = whitelist

	return nil
}
