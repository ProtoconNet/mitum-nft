package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *DelegateItem) unpack(
	enc encoder.Encoder,
	bAgent base.AddressDecoder,
	mode string,
	cid string,
) error {
	agent, err := bAgent.Encode(enc)
	if err != nil {
		return err
	}

	it.agent = agent
	it.mode = DelegateMode(mode)
	it.cid = currency.CurrencyID(cid)

	return nil
}
