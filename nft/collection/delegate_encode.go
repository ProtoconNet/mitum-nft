package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *DelegateFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bAgents []base.AddressDecoder,
	mode string,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	agents := make([]base.Address, len(bAgents))
	for i := range agents {
		agent, err := bAgents[i].Encode(enc)
		if err != nil {
			return err
		}

		agents[i] = agent
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.agents = agents
	fact.mode = DelegateMode(mode)
	fact.cid = currency.CurrencyID(cid)

	return nil
}
