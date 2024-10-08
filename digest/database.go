package digest

import (
	"context"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var maxLimit int64 = 50

var (
	defaultColNameAccount         = "digest_ac"
	defaultColNameContractAccount = "digest_ca"
	defaultColNameBalance         = "digest_bl"
	defaultColNameCurrency        = "digest_cr"
	defaultColNameOperation       = "digest_op"
	defaultColNameBlock           = "digest_bm"
	defaultColNameNFTCollection   = "digest_nftcollection"
	defaultColNameNFT             = "digest_nft"
	defaultColNameNFTOperator     = "digest_nftoperator"
)

func NFTCollection(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta mitumbase.State
	var err error
	if err := st.MongoClient().GetByFilter(
		defaultColNameNFTCollection,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			design, err = state.StateCollectionValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, mitumutil.ErrNotFound.WithMessage(err, "nft collection for contract account %v", contract)
	}

	return design, nil
}

func NFT(st *currencydigest.Database, contract, idx string) (*types.NFT, error) {
	i, err := strconv.ParseUint(idx, 10, 64)
	if err != nil {
		return nil, err
	}

	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("nft_idx", i)

	var nft *types.NFT
	var sta mitumbase.State
	if err = st.MongoClient().GetByFilter(
		defaultColNameNFT,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			nft, err = state.StateNFTValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, mitumutil.ErrNotFound.Errorf("nft token for contract account %v, nft idx %v", contract, idx)
	}

	return nft, nil
}

func NFTsByCollection(
	st *currencydigest.Database,
	contract, factHash, offset string,
	reverse bool,
	limit int64,
	callback func(nft types.NFT, st mitumbase.State) (bool, error),
) error {
	filter, err := buildNFTsFilterByContract(contract, factHash, offset, reverse)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("nft_idx", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.MongoClient().Find(
		context.Background(),
		defaultColNameNFT,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.Encoders())
			if err != nil {
				return false, err
			}
			nft, err := state.StateNFTValue(st)
			if err != nil {
				return false, err
			}
			return callback(*nft, st)
		},
		opt,
	)
}

func NFTCountByCollection(
	st *currencydigest.Database,
	contract string,
) (int64, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterToken := bson.D{{"istoken", true}}
	filterA = append(filterA, filterToken)
	filterA = append(filterA, filterContract)

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	opt := options.Count()

	return st.MongoClient().Count(
		context.Background(),
		defaultColNameNFT,
		filter,
		opt,
	)
}

func NFTOperators(
	st *currencydigest.Database,
	contract, account string,
) (*types.AllApprovedBook, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var operators *types.AllApprovedBook
	var sta mitumbase.State
	var err error
	if err := st.MongoClient().GetByFilter(
		defaultColNameNFTOperator,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			operators, err = state.StateOperatorsBookValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, mitumutil.ErrNotFound.WithMessage(err, "nft operators by contract %s and account %s", contract, account)
	}

	return operators, nil
}
