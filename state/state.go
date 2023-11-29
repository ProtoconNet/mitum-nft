package state

import (
	"github.com/ProtoconNet/mitum-nft/v2/types"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var CollectionStateValueHint = hint.MustNewHint("collection-state-value-v0.0.1")

type CollectionStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewCollectionStateValue(design types.Design) CollectionStateValue {
	return CollectionStateValue{
		BaseHinter: hint.NewBaseHinter(CollectionStateValueHint),
		Design:     design,
	}
}

func (cs CollectionStateValue) Hint() hint.Hint {
	return cs.BaseHinter.Hint()
}

func (cs CollectionStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionStateValue")

	if err := cs.BaseHinter.IsValid(CollectionStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := cs.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (cs CollectionStateValue) HashBytes() []byte {
	return cs.Design.Bytes()
}

func StateCollectionValue(st mitumbase.State) (*types.Design, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("collection design not found in State")
	}

	d, ok := v.(CollectionStateValue)
	if !ok {
		return nil, errors.Errorf("invalid collection value found, %T", v)
	}

	return &d.Design, nil
}

var LastNFTIndexStateValueHint = hint.MustNewHint("collection-last-nft-index-state-value-v0.0.1")

type LastNFTIndexStateValue struct {
	hint.BaseHinter
	id uint64
}

func NewLastNFTIndexStateValue( /*collection currencytypes.ContractID,*/ id uint64) LastNFTIndexStateValue {
	return LastNFTIndexStateValue{
		BaseHinter: hint.NewBaseHinter(LastNFTIndexStateValueHint),
		id:         id,
	}
}

func (is LastNFTIndexStateValue) Hint() hint.Hint {
	return is.BaseHinter.Hint()
}

func (is LastNFTIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid LastNFTIndexStateValue")

	if err := is.BaseHinter.IsValid(LastNFTIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (is LastNFTIndexStateValue) HashBytes() []byte {
	return util.Uint64ToBytes(is.id)
}

func StateLastNFTIndexValue(st mitumbase.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(LastNFTIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.id, nil
}

var (
	NFTStateValueHint = hint.MustNewHint("nft-state-value-v0.0.1")
)

type NFTStateValue struct {
	hint.BaseHinter
	NFT types.NFT
}

func NewNFTStateValue(n types.NFT) NFTStateValue {
	return NFTStateValue{
		BaseHinter: hint.NewBaseHinter(NFTStateValueHint),
		NFT:        n,
	}
}

func (ns NFTStateValue) Hint() hint.Hint {
	return ns.BaseHinter.Hint()
}

func (ns NFTStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NFTStateValue")

	if err := ns.BaseHinter.IsValid(NFTStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ns.NFT.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ns NFTStateValue) HashBytes() []byte {
	return ns.NFT.Bytes()
}

func StateNFTValue(st mitumbase.State) (*types.NFT, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("nft not found in State")
	}

	ns, ok := v.(NFTStateValue)
	if !ok {
		return nil, errors.Errorf("invalid nft value found, %T", v)
	}

	return &ns.NFT, nil
}

var NFTBoxStateValueHint = hint.MustNewHint("nft-box-state-value-v0.0.1")

type NFTBoxStateValue struct {
	hint.BaseHinter
	Box types.NFTBox
}

func NewNFTBoxStateValue(box types.NFTBox) NFTBoxStateValue {
	return NFTBoxStateValue{
		BaseHinter: hint.NewBaseHinter(NFTBoxStateValueHint),
		Box:        box,
	}
}

func (nb NFTBoxStateValue) Hint() hint.Hint {
	return nb.BaseHinter.Hint()
}

func (nb NFTBoxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid NFTBoxStateValue")

	if err := nb.BaseHinter.IsValid(NFTBoxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := nb.Box.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (nb NFTBoxStateValue) HashBytes() []byte {
	return nb.Box.Bytes()
}

func StateNFTBoxValue(st mitumbase.State) (types.NFTBox, error) {
	v := st.Value()
	if v == nil {
		return types.NFTBox{}, util.ErrNotFound.Errorf("nft box not found in State")
	}

	nb, ok := v.(NFTBoxStateValue)
	if !ok {
		return types.NFTBox{}, errors.Errorf("invalid nft box value found, %T", v)
	}

	return nb.Box, nil
}

var OperatorsBookStateValueHint = hint.MustNewHint("operators-book-state-value-v0.0.1")

type OperatorsBookStateValue struct {
	hint.BaseHinter
	Operators types.OperatorsBook
}

func NewOperatorsBookStateValue(operators types.OperatorsBook) OperatorsBookStateValue {
	return OperatorsBookStateValue{
		BaseHinter: hint.NewBaseHinter(OperatorsBookStateValueHint),
		Operators:  operators,
	}
}

func (ob OperatorsBookStateValue) Hint() hint.Hint {
	return ob.BaseHinter.Hint()
}

func (ob OperatorsBookStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid OperatorsBookStateValue")

	if err := ob.BaseHinter.IsValid(OperatorsBookStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ob.Operators.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ob OperatorsBookStateValue) HashBytes() []byte {
	return ob.Operators.Bytes()
}

func StateOperatorsBookValue(st mitumbase.State) (*types.OperatorsBook, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("operators book not found in State")
	}

	ob, ok := v.(OperatorsBookStateValue)
	if !ok {
		return nil, errors.Errorf("invalid operators book value found, %T", v)
	}

	return &ob.Operators, nil
}
