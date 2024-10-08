package digest

import (
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	crcystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type NFTCollectionDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de types.Design
}

func NewNFTCollectionDoc(st base.State, enc encoder.Encoder) (NFTCollectionDoc, error) {
	de, err := state.StateCollectionValue(st)
	if err != nil {
		return NFTCollectionDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return NFTCollectionDoc{}, err
	}

	return NFTCollectionDoc{
		BaseDoc: b,
		st:      st,
		de:      *de,
	}, nil
}

func (doc NFTCollectionDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["contract"] = doc.de.Contract()
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type NFTDoc struct {
	mongodbstorage.BaseDoc
	st        base.State
	nft       types.NFT
	addresses []base.Address
	owner     string
}

func NewNFTDoc(st base.State, enc encoder.Encoder) (*NFTDoc, error) {
	nft, err := state.StateNFTValue(st)
	if err != nil {
		return nil, err
	}
	var addresses = make([]string, len(nft.Creators().Addresses())+1)
	addresses[0] = nft.Owner().String()
	for i := range nft.Creators().Addresses() {
		addresses[i+1] = nft.Creators().Addresses()[i].String()
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTDoc{
		BaseDoc:   b,
		st:        st,
		nft:       *nft,
		addresses: nft.Addresses(),
		owner:     nft.Owner().String(),
	}, nil
}

func (doc NFTDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := crcystate.ParseStateKey(doc.st.Key(), state.NFTPrefix, 4)
	if err != nil {
		return nil, err
	}

	var hashArray []string
	for _, v := range doc.st.Operations() {
		hashArray = append(hashArray, v.String())
	}

	m["contract"] = parsedKey[1]
	m["nft_idx"] = doc.nft.ID()
	m["owner"] = doc.nft.Owner()
	m["addresses"] = doc.addresses
	m["istoken"] = true
	m["height"] = doc.st.Height()
	m["facthash"] = hashArray

	return bsonenc.Marshal(m)
}

type NFTAllApprovedDoc struct {
	mongodbstorage.BaseDoc
	st        base.State
	operators types.AllApprovedBook
}

func NewNFTOperatorDoc(st base.State, enc encoder.Encoder) (*NFTAllApprovedDoc, error) {
	operators, err := state.StateOperatorsBookValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTAllApprovedDoc{
		BaseDoc:   b,
		st:        st,
		operators: *operators,
	}, nil
}

func (doc NFTAllApprovedDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := crcystate.ParseStateKey(doc.st.Key(), state.NFTPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["address"] = parsedKey[2]
	m["approved"] = doc.operators
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type NFTLastIndexDoc struct {
	mongodbstorage.BaseDoc
	st    base.State
	nftID uint64
}

func NewNFTLastIndexDoc(st base.State, enc encoder.Encoder) (*NFTLastIndexDoc, error) {
	nftID, err := state.StateLastNFTIndexValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTLastIndexDoc{
		BaseDoc: b,
		st:      st,
		nftID:   nftID,
	}, nil
}

func (doc NFTLastIndexDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := crcystate.ParseStateKey(doc.st.Key(), state.NFTPrefix, 3)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["nft_idx"] = doc.nftID
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
