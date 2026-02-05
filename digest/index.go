package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var nftCollectionIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "nft_collection_contract_height"),
	},
}

var nftIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "nft_idx", Value: 1},
			bson.E{Key: "height", Value: -1},
			bson.E{Key: "istoken", Value: 1},
		},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "nft_idx_contract_height_istoken"),
	},
	{
		Keys: bson.D{bson.E{Key: "facthash", Value: 1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "nft_facthash"),
	},
}

var nftOperatorIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "nft_operator_contract_address_height"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNameNFTCollection] = nftCollectionIndexModels
	DefaultIndexes[DefaultColNameNFT] = nftIndexModels
	DefaultIndexes[DefaultColNameNFTOperator] = nftOperatorIndexModels
}
