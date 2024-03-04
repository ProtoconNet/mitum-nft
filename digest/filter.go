package digest

import (
	"net/http"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func parseRequest(w http.ResponseWriter, r *http.Request, v string) (string, error, int) {
	s, found := mux.Vars(r)[v]
	if !found {
		return "", errors.Errorf("empty %s", v), http.StatusNotFound
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return "", errors.Errorf("empty %s", v), http.StatusBadRequest
	}
	return s, nil, http.StatusOK
}

func buildNFTsFilterByAddress(address base.Address, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching address
	filterAddress := bson.D{{"owner", bson.D{{"$in", []string{address.String()}}}}}
	filterA = append(filterA, filterAddress)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"nftid", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
			// if reverse true, lesser then offset height
		} else {
			filterHeight := bson.D{
				{"nftid", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
		}
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}

func buildNFTsFilterByContract(contract, offset, facthash string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterToken := bson.D{{"istoken", true}}
	filterA = append(filterA, filterToken)
	filterA = append(filterA, filterContract)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"nftid", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
			// if reverse true, lesser then offset height
		} else {
			filterHeight := bson.D{
				{"nftid", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
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
