package cmds

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

type KeyFlag struct {
	Key currencybase.BaseAccountKey
}

func (v *KeyFlag) UnmarshalText(b []byte) error {
	if bytes.Equal(bytes.TrimSpace(b), []byte("-")) {
		c, err := LoadFromStdInput()
		if err != nil {
			return err
		}
		b = c
	}

	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return errors.Errorf(`wrong formatted; "<string private key>,<uint weight>"`)
	}

	var pk base.Publickey
	if k, err := base.DecodePublickeyFromString(l[0], enc); err != nil {
		return errors.Wrapf(err, "invalid public key, %q for --key", l[0])
	} else {
		pk = k
	}

	var weight uint = 100
	if i, err := strconv.ParseUint(l[1], 10, 8); err != nil {
		return errors.Wrapf(err, "invalid weight, %q for --key", l[1])
	} else if i > 0 && i <= 100 {
		weight = uint(i)
	}

	if k, err := currencybase.NewBaseAccountKey(pk, weight); err != nil {
		return err
	} else if err := k.IsValid(nil); err != nil {
		return errors.Wrap(err, "invalid key string")
	} else {
		v.Key = k
	}

	return nil
}

type StringLoad []byte

func (v *StringLoad) UnmarshalText(b []byte) error {
	if bytes.Equal(bytes.TrimSpace(b), []byte("-")) {
		c, err := LoadFromStdInput()
		if err != nil {
			return err
		}
		*v = c

		return nil
	}

	*v = b

	return nil
}

func (v StringLoad) Bytes() []byte {
	return []byte(v)
}

func (v StringLoad) String() string {
	return string(v)
}

type PrivatekeyFlag struct {
	base.Privatekey
	notEmpty bool
}

func (v PrivatekeyFlag) Empty() bool {
	return !v.notEmpty
}

func (v *PrivatekeyFlag) UnmarshalText(b []byte) error {
	if k, err := base.DecodePrivatekeyFromString(string(b), enc); err != nil {
		return errors.Wrapf(err, "invalid private key, %q", string(b))
	} else if err := k.IsValid(nil); err != nil {
		return err
	} else {
		*v = PrivatekeyFlag{Privatekey: k}
	}

	v.notEmpty = true

	return nil
}

type PublickeyFlag struct {
	base.Publickey
	notEmpty bool
}

func (v PublickeyFlag) Empty() bool {
	return !v.notEmpty
}

func (v *PublickeyFlag) UnmarshalText(b []byte) error {
	if k, err := base.DecodePublickeyFromString(string(b), enc); err != nil {
		return errors.Wrapf(err, "invalid public key, %q", string(b))
	} else if err := k.IsValid(nil); err != nil {
		return err
	} else {
		*v = PublickeyFlag{Publickey: k}
	}

	v.notEmpty = true

	return nil
}

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
	return base.DecodeAddress(v.s, enc)
}

type BigFlag struct {
	currencybase.Big
}

func (v *BigFlag) UnmarshalText(b []byte) error {
	if a, err := currencybase.NewBigFromString(string(b)); err != nil {
		return errors.Wrapf(err, "invalid big string, %q", string(b))
	} else if err := a.IsValid(nil); err != nil {
		return err
	} else {
		*v = BigFlag{Big: a}
	}

	return nil
}

type CurrencyIDFlag struct {
	CID currencybase.CurrencyID
}

func (v *CurrencyIDFlag) UnmarshalText(b []byte) error {
	cid := currencybase.CurrencyID(string(b))
	if err := cid.IsValid(nil); err != nil {
		return err
	}
	v.CID = cid

	return nil
}

func (v *CurrencyIDFlag) String() string {
	return v.CID.String()
}

type ContractIDFlag struct {
	ID currencybase.ContractID
}

func (v *ContractIDFlag) UnmarshalText(b []byte) error {
	id := currencybase.ContractID(string(b))
	if err := id.IsValid(nil); err != nil {
		return err
	}
	v.ID = id

	return nil
}

func (v *ContractIDFlag) String() string {
	return v.ID.String()
}

type CurrencyAmountFlag struct {
	CID currencybase.CurrencyID
	Big currencybase.Big
}

func (v *CurrencyAmountFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid currency-amount, %q", string(b))
	}

	a, c := l[0], l[1]

	cid := currencybase.CurrencyID(a)
	if err := cid.IsValid(nil); err != nil {
		return err
	}
	v.CID = cid

	if a, err := currencybase.NewBigFromString(c); err != nil {
		return errors.Wrapf(err, "invalid big string, %q", string(b))
	} else if err := a.IsValid(nil); err != nil {
		return err
	} else {
		v.Big = a
	}

	return nil
}

func (v *CurrencyAmountFlag) String() string {
	return v.CID.String() + "," + v.Big.String()
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
	collection currencybase.ContractID
	idx        uint64
}

func (v *NFTIDFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), "-", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid nft id, %q", string(b))
	}

	s, id := l[0], l[1]

	collection := currencybase.ContractID(s)
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
