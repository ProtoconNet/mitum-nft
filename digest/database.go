package digest

import (
	"context"
	"strconv"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	cutil "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var maxLimit int64 = 50

var (
	DefaultColNameNFTCollection = "digest_nftcollection"
	DefaultColNameNFT           = "digest_nft"
	DefaultColNameNFTOperator   = "digest_nftoperator"
)

func NFTCollection(st *cdigest.Database, contract string) (*types.Design, error) {
	filter := cutil.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNameNFTCollection,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			design, err = state.StateCollectionValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(cutil.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, util.ErrNotFound.WithMessage(err, "nft collection for contract account %v", contract)
	}

	return design, nil
}

func NFT(st *cdigest.Database, contract, idx string) (*types.NFT, error) {
	i, err := strconv.ParseUint(idx, 10, 64)
	if err != nil {
		return nil, err
	}

	filter := cutil.NewBSONFilter("contract", contract)
	filter = filter.Add("nft_idx", i)

	var nft *types.NFT
	var sta base.State
	if err = st.MongoClient().GetByFilter(
		DefaultColNameNFT,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			nft, err = state.StateNFTValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(cutil.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, util.ErrNotFound.Errorf("nft token for contract account %v, nft idx %v", contract, idx)
	}

	return nft, nil
}

func NFTsByCollection(
	st *cdigest.Database,
	contract, factHash, offset string,
	reverse bool,
	limit int64,
	callback func(nft types.NFT, st base.State) (bool, error),
) error {
	sortDir := 1
	cmpOp := "$gt"
	if reverse {
		sortDir = -1
		cmpOp = "$lt"
	}

	match := bson.D{
		{Key: "contract", Value: contract},
	}

	if factHash != "" {
		match = append(match, bson.E{Key: "facthash", Value: factHash})
	}

	if offset != "" {
		match = append(match, bson.E{
			Key:   "nft_idx",
			Value: bson.D{{Key: cmpOp, Value: offset}},
		})
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "nft_idx", Value: sortDir},
			{Key: "height", Value: -1},
			{Key: "_id", Value: -1},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$nft_idx"},
			{Key: "doc", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id", Value: sortDir},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: "$doc"},
		}}},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	}

	return st.MongoClient().Aggregate(
		context.Background(),
		DefaultColNameNFT,
		pipeline,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := cdigest.LoadState(cursor.Decode, st.Encoders())
			if err != nil {
				return false, err
			}
			nft, err := state.StateNFTValue(st)
			if err != nil {
				return false, err
			}
			return callback(*nft, st)
		},
	)
}

func NFTOperators(
	st *cdigest.Database,
	contract, account string,
) (*types.AllApprovedBook, error) {
	filter := cutil.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var operators *types.AllApprovedBook
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNameNFTOperator,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			operators, err = state.StateOperatorsBookValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(cutil.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, util.ErrNotFound.WithMessage(err, "nft operators by contract %s and account %s", contract, account)
	}

	return operators, nil
}
