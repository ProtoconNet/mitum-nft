package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (d *Design) unpack(
	enc encoder.Encoder,
	bpr base.AddressDecoder,
	bc base.AddressDecoder,
	symbol string,
	active bool,
	bpo []byte,
) error {

	parent, err := bpr.Encode(enc)
	if err != nil {
		return err
	}
	d.parent = parent

	creator, err := bc.Encode(enc)
	if err != nil {
		return err
	}
	d.creator = creator

	d.symbol = extensioncurrency.ContractID(symbol)
	d.active = active

	if hinter, err := enc.Decode(bpo); err != nil {
		return err
	} else if policy, ok := hinter.(BasePolicy); !ok {
		return util.WrongTypeError.Errorf("not BasePolicy; %T", hinter)
	} else {
		d.policy = policy
	}

	return nil
}
