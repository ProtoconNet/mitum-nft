package collection

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CollectionRegisterFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	ca string,
	sb string,
	nm string,
	ry uint,
	uri string,
	bws []string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionRegisterFact")

	fact.currency = currencybase.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	fact.collection = currencybase.ContractID(sb)
	fact.name = CollectionName(nm)
	fact.royalty = nft.PaymentParameter(ry)
	fact.uri = nft.URI(uri)

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return e(err, "")
	}
	fact.contract = contract

	whites := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e(err, "")
		}
		whites[i] = white

	}
	fact.whitelist = whites

	return nil
}
