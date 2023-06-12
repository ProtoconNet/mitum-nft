package digest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func parseRequest(w http.ResponseWriter, r *http.Request, v string) (string, error, int) {
	s, found := mux.Vars(r)[v]
	if !found {
		return "", errors.Errorf("empty collection id"), http.StatusNotFound
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return "", errors.Errorf("empty collection id"), http.StatusBadRequest
	}
	return s, nil, http.StatusOK
}

func parseOffset(s string) (base.Height, uint64, error) {
	if n := strings.SplitN(s, ",", 2); n == nil {
		return base.NilHeight, 0, errors.Errorf("invalid offset string: %q", s)
	} else if len(n) < 2 {
		return base.NilHeight, 0, errors.Errorf("invalid offset, %q", s)
	} else if h, err := base.ParseHeightString(n[0]); err != nil {
		return base.NilHeight, 0, errors.Wrap(err, "invalid height of offset")
	} else if u, err := strconv.ParseUint(n[1], 10, 64); err != nil {
		return base.NilHeight, 0, errors.Wrap(err, "invalid index of offset")
	} else {
		return h, u, nil
	}
}

func buildOffset(height base.Height, index uint64) string {
	return fmt.Sprintf("%d,%d", height, index)
}

func buildOperationsFilterByAddress(address base.Address, offset string, reverse bool) (bson.M, error) {
	filter := bson.M{"addresses": bson.M{"$in": []string{address.String()}}}
	if len(offset) > 0 {
		height, index, err := parseOffset(offset)
		if err != nil {
			return nil, err
		}

		if reverse {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$lt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"index": bson.M{"$lt": index}},
				}},
			}
		} else {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$gt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"index": bson.M{"$gt": index}},
				}},
			}
		}
	}

	return filter, nil
}

func parseOffsetByString(s string) (base.Height, string, error) {
	var a, b string
	switch n := strings.SplitN(s, ",", 2); {
	case n == nil:
		return base.NilHeight, "", errors.Errorf("invalid offset string: %q", s)
	case len(n) < 2:
		return base.NilHeight, "", errors.Errorf("invalid offset, %q", s)
	default:
		a = n[0]
		b = n[1]
	}

	h, err := base.ParseHeightString(a)
	if err != nil {
		return base.NilHeight, "", errors.Wrap(err, "invalid height of offset")
	}

	return h, b, nil
}

func buildOffsetByString(height base.Height, s string) string {
	return fmt.Sprintf("%d,%s", height, s)
}

func buildAccountsFilterByPublickey(pub base.Publickey) bson.M {
	return bson.M{"pubs": bson.M{"$in": []string{pub.String()}}}
}

func buildNFTsFilterByAddress(address base.Address, offset string, reverse bool, collection string) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching address
	filterAddress := bson.D{{"owner", bson.D{{"$in", []string{address.String()}}}}}
	filterA = append(filterA, filterAddress)

	// if collection query exist, find by collection first
	if len(collection) > 0 {
		filterCollection := bson.D{
			{"collection", bson.D{{"$eq", collection}}},
		}
		filterA = append(filterA, filterCollection)
	}

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

func buildNFTsFilterByCollection(contract, col string, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterSymbol := bson.D{{"collection", bson.D{{"$in", []string{col}}}}}
	filterToken := bson.D{{"istoken", true}}
	filterA = append(filterA, filterToken)
	filterA = append(filterA, filterSymbol)

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
