package digest

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func buildNFTsFilterByContract(contract, facthash, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterToken := bson.D{{"istoken", true}}
	filterA = append(filterA, filterToken)
	filterA = append(filterA, filterContract)

	// if offset exist, apply offset
	if len(offset) > 0 {
		v, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return nil, err
		}

		if !reverse {
			filterOffset := bson.D{
				{"nft_idx", bson.D{{"$gt", v}}},
			}
			filterA = append(filterA, filterOffset)
		} else {
			filterOffset := bson.D{
				{"nft_idx", bson.D{{"$lt", v}}},
			}
			filterA = append(filterA, filterOffset)
		}
	}

	if len(facthash) > 0 {
		filterFactHash := bson.D{
			{"facthash", bson.D{{"$in", []string{facthash}}}},
		}
		filterA = append(filterA, filterFactHash)
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}
