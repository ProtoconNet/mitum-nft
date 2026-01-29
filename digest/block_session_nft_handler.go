package digest

import (
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareNFTs() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var nftCollectionModels []mongo.WriteModel
	var nftOperatorModels []mongo.WriteModel
	var nftModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		stateKey, err := state.ParseNFTStateKey(st.Key())
		if err != nil {
			continue
		}
		switch stateKey {
		case state.CollectionKey:
			j, err := bs.handleNFTCollectionState(st)
			if err != nil {
				return err
			}
			nftCollectionModels = append(nftCollectionModels, j...)
		case state.OperatorsKey:
			j, err := bs.handleNFTOperatorsState(st)
			if err != nil {
				return err
			}
			nftOperatorModels = append(nftOperatorModels, j...)
		case state.NFTKey:
			j, err := bs.handleNFTState(st)
			if err != nil {
				return err
			}
			nftModels = append(nftModels, j...)
			bs.nftMap[st.Key()] = struct{}{}
		default:
			continue
		}
	}

	bs.nftCollectionModels = nftCollectionModels
	bs.nftOperatorModels = nftOperatorModels
	bs.nftModels = nftModels

	return nil
}

func (bs *BlockSession) handleNFTCollectionState(st base.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewNFTCollectionDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func (bs *BlockSession) handleNFTOperatorsState(st base.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewNFTOperatorDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func (bs *BlockSession) handleNFTState(st base.State) ([]mongo.WriteModel, error) {
	if nftDoc, err := NewNFTDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftDoc),
		}, nil
	}
}

func (bs *BlockSession) handleNFTLastIndexState(st base.State) ([]mongo.WriteModel, error) {
	if nftLastIndexDoc, err := NewNFTLastIndexDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftLastIndexDoc),
		}, nil
	}
}
