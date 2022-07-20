package digest

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

type NFTCollectionDoc struct {
	mongodbstorage.BaseDoc
	st state.State
	de nft.Design
}

func NewNFTCollectionDoc(st state.State, enc encoder.Encoder) (NFTCollectionDoc, error) {
	de, err := collection.StateCollectionValue(st)
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
		de:      de,
	}, nil
}

func (doc NFTCollectionDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["symbol"] = doc.de.Symbol()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type NFTDoc struct {
	mongodbstorage.BaseDoc
	va        NFTValue
	addresses []string
	owner     string
	height    base.Height
}

func NewNFTDoc(st state.State, enc encoder.Encoder, height base.Height) (NFTDoc, error) {
	n, err := collection.StateNFTValue(st)
	if err != nil {
		return NFTDoc{}, err
	}
	var addresses = make([]string, len(n.Creators().Addresses())+len(n.Copyrighters().Addresses())+1)
	addresses[0] = n.Owner().String()
	for i := range n.Creators().Addresses() {
		addresses[i+1] = n.Creators().Addresses()[i].String()
	}
	for i := range n.Copyrighters().Addresses() {
		addresses[i+1+len(n.Creators().Addresses())] = n.Copyrighters().Addresses()[i].String()
	}
	va := NewNFTValue(n, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return NFTDoc{}, err
	}

	return NFTDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		owner:     n.Owner().String(),
		height:    height,
	}, nil
}

func (doc NFTDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["collection"] = doc.va.nft.ID().Collection()
	m["nftid"] = doc.va.nft.ID().String()
	m["owner"] = doc.va.nft.Owner()
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}

type NFTAgentDoc struct {
	mongodbstorage.BaseDoc
	st     state.State
	agents collection.AgentBox
}

func NewNFTAgentDoc(st state.State, enc encoder.Encoder) (NFTAgentDoc, error) {
	agents, err := collection.StateAgentsValue(st)
	if err != nil {
		return NFTAgentDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return NFTAgentDoc{}, err
	}

	return NFTAgentDoc{
		BaseDoc: b,
		st:      st,
		agents:  agents,
	}, nil
}

func (doc NFTAgentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	address := doc.st.Key()[:len(doc.st.Key())-len(doc.agents.Collection().String())-len(collection.StateKeyAgentsSuffix)-1]
	m["collectionid"] = doc.agents.Collection().String()
	m["address"] = address
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
