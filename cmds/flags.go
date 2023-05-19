package cmds

import (
	"fmt"
	"strconv"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/v2/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

type ContractIDFlag struct {
	CID extensioncurrency.ContractID
}

func (v *ContractIDFlag) UnmarshalText(b []byte) error {
	cid := extensioncurrency.ContractID(string(b))
	if err := cid.IsValid(nil); err != nil {
		return err
	}
	v.CID = cid

	return nil
}

func (v *ContractIDFlag) String() string {
	return v.CID.String()
}

type SignerFlag struct {
	address string
	share   uint
}

func (v *SignerFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid signer; %q", string(b))
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
	return base.DecodeAddress(v.address, enc)
}

type NFTIDFlag struct {
	collection extensioncurrency.ContractID
	idx        uint64
}

func (v *NFTIDFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), "-", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid nft id, %q", string(b))
	}

	s, id := l[0], l[1]

	collection := extensioncurrency.ContractID(s)
	if err := collection.IsValid(nil); err != nil {
		return err
	}
	v.collection = collection

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
