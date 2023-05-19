package collection

import (
	"fmt"
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CollectionPolicyUpdaterFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	ct string,
	col string,
	nm string,
	ry uint,
	uri string,
	bws []byte,
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

	contract, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.contract = contract

	fact.name = CollectionName(nm)
	fact.royalty = nft.PaymentParameter(ry)
	fact.uri = nft.URI(uri)

	hits, err := enc.DecodeSlice(bws)
	if err != nil {
		return e(err, "")
	}

	whitelist := make([]base.Address, len(bws))
	for i := range hits {
		ad := fmt.Sprintf("%v", hits[i])
		white, err := base.DecodeAddress(ad, enc)
		if err != nil {
			return e(err, "")
		}
		whitelist[i] = white
	}
	fact.whitelist = whitelist
	//whitelist := make([]base.Address, len(bws))
	//for i, bw := range bws {
	//	white, err := base.DecodeAddress(bw, enc)
	//	if err != nil {
	//		return e(err, "")
	//	}
	//	whitelist[i] = white
	//}
	//fact.whitelist = whitelist

	return nil
}
