package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func PrepareNFTs(bs *currencydigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
	stateKey, err := state.ParseNFTStateKey(st.Key())
	if err != nil {
		return "", nil, nil
	}
	switch stateKey {
	case state.CollectionKey:
		j, err := handleNFTCollectionState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameNFTCollection, j, nil
	case state.OperatorsKey:
		j, err := handleNFTOperatorsState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameNFTOperator, j, nil
	case state.NFTKey:
		j, err := handleNFTState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameNFT, j, nil
	}

	return "", nil, nil
}

func handleNFTCollectionState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewNFTCollectionDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func handleNFTOperatorsState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewNFTOperatorDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func handleNFTState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if nftDoc, err := NewNFTDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftDoc),
		}, nil
	}
}
