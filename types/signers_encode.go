package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (sgns *Signers) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	bsns []byte,
) error {
	e := util.StringError("failed to unmarshal Signers")

	sgns.BaseHinter = hint.NewBaseHinter(ht)

	hinters, err := enc.DecodeSlice(bsns)
	if err != nil {
		return e.Wrap(err)
	}

	signers := make([]Signer, len(hinters))
	for i, hinter := range hinters {
		signer, ok := hinter.(Signer)
		if !ok {
			return e.Wrap(errors.Errorf("expected Signer, not %T", hinter))
		}

		signers[i] = signer
	}
	sgns.signers = signers

	return nil
}
