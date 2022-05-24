package nft

import (
	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (d *Design) unpack(
	enc encoder.Encoder,
	_symbol string,
	bPolicy []byte,
) error {
	d.symbol = extension.ContractID(_symbol)

	var policy BasePolicy
	if hinter, err := enc.Decode(bPolicy); err != nil {
		return err
	} else if i, ok := hinter.(BasePolicy); !ok {
		return util.WrongTypeError.Errorf("not BasePolicy; %T", hinter)
	} else {
		policy = i
	}
	d.policy = policy

	return nil
}
