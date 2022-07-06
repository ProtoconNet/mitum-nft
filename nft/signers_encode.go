package nft

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (signers *Signers) unpack(
	enc encoder.Encoder,
	total uint,
	bsns []byte,
) error {
	signers.total = total

	hsns, err := enc.DecodeSlice(bsns)
	if err != nil {
		return err
	}

	sns := make([]Signer, len(hsns))
	for i := range hsns {
		signer, ok := hsns[i].(Signer)
		if !ok {
			return util.WrongTypeError.Errorf("not Signer; %T", hsns[i])
		}
		sns[i] = signer
	}
	signers.signers = sns

	return nil
}
