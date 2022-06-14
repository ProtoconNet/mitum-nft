package cmds

import (
	"fmt"
	"strconv"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
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

type RightHolerFlag struct {
	address string
	clue    string
}

func (v *RightHolerFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid right holder; %q", string(b))
	}
	v.address, v.clue = l[0], l[1]

	return nil
}

func (v *RightHolerFlag) String() string {
	if len(v.address) > 0 || len(v.clue) > 0 {
		return ""
	}
	s := fmt.Sprintf("%s,%s", v.address, v.clue)
	return s
}

func (v *RightHolerFlag) Encode(enc encoder.Encoder) (nft.RightHoler, error) {
	account, err := base.DecodeAddressFromString(v.address, enc)
	if err != nil {
		return nft.RightHoler{}, err
	}

	r := nft.NewRightHoler(account, false, v.clue)
	if err != nil {
		return nft.RightHoler{}, err
	}

	return r, nil
}
