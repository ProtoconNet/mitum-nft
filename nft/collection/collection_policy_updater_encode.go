package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CollectionPolicyUpdaterFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	col string,
	nm string,
	ry uint,
	uri string,
	bws []string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionPolicyUpdaterFact")

	fact.collection = extensioncurrency.ContractID(col)
	fact.currency = currency.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	fact.name = CollectionName(nm)
	fact.royalty = nft.PaymentParameter(ry)
	fact.uri = nft.URI(uri)

	whitelist := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e(err, "")
		}
		whitelist[i] = white
	}
	fact.whitelist = whitelist

	return nil
}
