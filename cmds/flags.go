package cmds

import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

type AddressFlag struct {
	s string
}

func (v *AddressFlag) UnmarshalText(b []byte) error {
	v.s = string(b)

	return nil
}

func (v *AddressFlag) String() string {
	return v.s
}

func (v *AddressFlag) Encode(enc encoder.Encoder) (base.Address, error) {
	return base.DecodeAddressFromString(v.s, enc)
}

type NFTIDFlag struct {
	collection extension.ContractID
	idx        currency.Big
}

func (v *NFTIDFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid nft id; %q", string(b))
	}

	s, id := l[0], l[1]

	symbol := extension.ContractID(s)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}
	v.collection = symbol

	if idx, err := currency.NewBigFromString(id); err != nil {
		return errors.Wrapf(err, "invalid big string; %q", string(b))
	} else if err := idx.IsValid(nil); err != nil {
		return err
	} else {
		v.idx = idx
	}

	return nil
}

func (v *NFTIDFlag) String() string {
	return v.collection.String() + "," + v.idx.String()
}
