package cmds

import (
	"fmt"
	"strconv"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
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
	collection extensioncurrency.ContractID
	idx        uint64
}

func (v *NFTIDFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid nft id; %q", string(b))
	}

	s, id := l[0], l[1]

	symbol := extensioncurrency.ContractID(s)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}
	v.collection = symbol

	if i, err := strconv.ParseUint(id, 10, 64); err != nil {
		return err
	} else {
		v.idx = i
	}

	return nil
}

func (v *NFTIDFlag) String() string {
	s := fmt.Sprintf("%s,%d", v.collection, v.idx)
	return s
}

type SignerFlag struct {
	address string
	share   uint
}

func (v *SignerFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid nft id; %q", string(b))
	}

	v.address = l[0]

	if share, err := strconv.ParseUint(l[1], 10, 8); err != nil {
		return err
	} else if share > uint64(nft.MaxSignerShare) {
		return errors.Errorf("share is over max; %d > %d", share, nft.MaxSignerShare)
	} else {
		v.share = uint(share)
	}

	return nil
}

func (v *SignerFlag) String() string {
	s := fmt.Sprintf("%s,%d", v.address, v.share)
	return s
}

func (v *SignerFlag) Encode(enc encoder.Encoder) (base.Address, error) {
	return base.DecodeAddressFromString(v.address, enc)
}
