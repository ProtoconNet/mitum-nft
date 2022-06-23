package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *DelegateItem) unpack(
	enc encoder.Encoder,
	collection string,
	bag base.AddressDecoder,
	mode string,
	cid string,
) error {
	it.collection = extensioncurrency.ContractID(collection)

	agent, err := bag.Encode(enc)
	if err != nil {
		return err
	}
	it.agent = agent

	it.mode = DelegateMode(mode)
	it.cid = currency.CurrencyID(cid)

	return nil
}
