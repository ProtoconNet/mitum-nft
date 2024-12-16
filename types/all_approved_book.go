package types

import (
	"bytes"
	"sort"

	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxAllApproved = 10

var AllApprovedBookHint = hint.MustNewHint("mitum-nft-all-approved-book-v0.0.1")

type AllApprovedBook struct {
	hint.BaseHinter
	allApproved []base.Address
}

func NewAllApprovedBook(allApproved []base.Address) AllApprovedBook {
	if allApproved == nil {
		return AllApprovedBook{
			BaseHinter:  hint.NewBaseHinter(AllApprovedBookHint),
			allApproved: []base.Address{},
		}
	}
	return AllApprovedBook{
		BaseHinter:  hint.NewBaseHinter(AllApprovedBookHint),
		allApproved: allApproved,
	}
}

func (ob AllApprovedBook) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		ob.BaseHinter,
	); err != nil {
		return err
	}
	for i := range ob.allApproved {
		if err := ob.allApproved[i].IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (ob AllApprovedBook) Bytes() []byte {
	aaps := make([][]byte, len(ob.allApproved))

	for i, aap := range ob.allApproved {
		aaps[i] = aap.Bytes()
	}

	return util.ConcatBytesSlice(aaps...)
}

func (ob AllApprovedBook) Hash() util.Hash {
	return ob.GenerateHash()
}

func (ob AllApprovedBook) GenerateHash() util.Hash {
	return valuehash.NewSHA256(ob.Bytes())
}

func (ob AllApprovedBook) IsEmpty() bool {
	return len(ob.allApproved) < 1
}

func (ob AllApprovedBook) Equal(b AllApprovedBook) bool {
	ob.Sort(true)
	b.Sort(true)

	for i := range ob.allApproved {
		if !ob.allApproved[i].Equal(b.allApproved[i]) {
			return false
		}
	}

	return true
}

func (ob *AllApprovedBook) Sort(ascending bool) {
	sort.Slice(ob.allApproved, func(i, j int) bool {
		if ascending {
			return bytes.Compare(ob.allApproved[j].Bytes(), ob.allApproved[i].Bytes()) > 0
		}

		return bytes.Compare(ob.allApproved[j].Bytes(), ob.allApproved[i].Bytes()) < 0
	})
}

func (ob AllApprovedBook) Exists(ag base.Address) bool {
	if ob.IsEmpty() {
		return false
	}

	for _, operator := range ob.allApproved {
		if ag.Equal(operator) {
			return true
		}
	}

	return false
}

func (ob AllApprovedBook) Get(ag base.Address) (base.Address, error) {
	for _, operator := range ob.allApproved {
		if ag.Equal(operator) {
			return operator, nil
		}
	}

	return ctypes.Address{}, errors.Errorf("account %v not in operators book", ag)
}

func (ob *AllApprovedBook) Append(ag base.Address) error {
	if err := ag.IsValid(nil); err != nil {
		return err
	}

	if ob.Exists(ag) {
		return errors.Errorf("account already in operators book, %v", ag)
	}

	if len(ob.allApproved) >= MaxAllApproved {
		return errors.Errorf("max operators, %v", ag)
	}

	ob.allApproved = append(ob.allApproved, ag)

	return nil
}

func (ob *AllApprovedBook) Remove(ag base.Address) error {
	if !ob.Exists(ag) {
		return errors.Errorf("account %v not in operators book", ag)
	}

	for i := range ob.allApproved {
		if ag.String() == ob.allApproved[i].String() {
			ob.allApproved[i] = ob.allApproved[len(ob.allApproved)-1]
			ob.allApproved[len(ob.allApproved)-1] = ctypes.Address{}
			ob.allApproved = ob.allApproved[:len(ob.allApproved)-1]

			return nil
		}
	}
	return nil
}

func (ob AllApprovedBook) AllApproved() []base.Address {
	return ob.allApproved
}
