package digest

import (
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleNFT(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	id, err, status := parseRequest(w, r, "id")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTInGroup(contract, collection, id)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)
		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTInGroup(contract, collection, id string) (interface{}, error) {
	switch nft, st, err := hd.database.NFT(contract, collection, id); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTHal(contract, collection, *nft, st)
		if err != nil {
			return nil, err
		}
		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTHal(contract, collection string, nft types.NFT, st base.State) (Hal, error) {
	h, err := hd.combineURL(HandlerPathNFT, "contract", contract, "collection", collection, "id", string(nft.ID()))
	if err != nil {
		return nil, err
	}

	hal := NewBaseHal(nft, NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleNFTCollection(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTCollectionInGroup(contract, collection)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)
		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTCollectionInGroup(contract, collection string) (interface{}, error) {
	switch design, st, err := hd.database.NFTCollection(contract, collection); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTCollectionHal(contract, collection, *design, st)
		if err != nil {
			return nil, err
		}
		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTCollectionHal(contract, collection string, design types.Design, st base.State) (Hal, error) {
	h, err := hd.combineURL(HandlerPathNFTs, "contract", contract, "collection", collection)
	if err != nil {
		return nil, err
	}

	hal := NewBaseHal(design, NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleNFTs(w http.ResponseWriter, r *http.Request) {
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseStringQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := CacheKey(
		r.URL.Path, stringOffsetQuery(offset),
		stringBoolQuery("reverse", reverse),
	)

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	// if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
	// 	return hd.handleCollectionNFTsInGroup(contract, collection)
	// }); err != nil {
	// 	HTTP2HandleError(w, err)
	// } else {
	// 	HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)

	// 	if !shared {
	// 		HTTP2WriteCache(w, cachekey, time.Second*3)
	// 	}
	// }

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleNFTsInGroup(contract, collection, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("collection", collection).Msg("failed to get nfts")
		HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	HTTP2WriteHalBytes(hd.enc, w, b, http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		HTTP2WriteCache(w, cachekey, expire)
	}
}

func (hd *Handlers) handleNFTsInGroup(
	contract, collection string,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("collection-nfts")
	} else {
		limit = l
	}

	var vas []Hal
	if err := hd.database.NFTsByCollection(
		contract, collection, reverse, offset, limit,
		func(nft types.NFT, st base.State) (bool, error) {
			hal, err := hd.buildNFTHal(contract, collection, nft, st)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, err
	} else if len(vas) < 1 {
		return nil, false, errors.Errorf("nfts not found")
	}

	i, err := hd.buildCollectionNFTsHal(contract, collection, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.enc.Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) buildCollectionNFTsHal(
	contract, col string,
	vas []Hal,
	offset string,
	reverse bool,
) (Hal, error) {
	baseSelf, err := hd.combineURL(HandlerPathNFTs, "contract", contract, "collection", col)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = addQueryValue(baseSelf, stringOffsetQuery(offset))
	}
	if reverse {
		self = addQueryValue(baseSelf, stringBoolQuery("reverse", reverse))
	}

	var hal Hal
	hal = NewBaseHal(vas, NewHalLink(self, nil))

	h, err := hd.combineURL(HandlerPathNFTCollection, "contract", contract, "collection", col)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("collection", NewHalLink(h, nil))

	var nextoffset string

	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(types.NFT)
		nextoffset = string(va.ID())
	}

	if len(nextoffset) > 0 {
		next := baseSelf
		next = addQueryValue(next, stringOffsetQuery(nextoffset))

		if reverse {
			next = addQueryValue(next, stringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", NewHalLink(next, nil))
	}

	hal = hal.AddLink("reverse", NewHalLink(addQueryValue(baseSelf, stringBoolQuery("reverse", !reverse)), nil))

	return hal, nil
}

func (hd *Handlers) handleNFTOperators(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := parseRequest(w, r, "account")
	if err != nil {
		HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTOperatorsInGroup(contract, collection, account)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)
		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTOperatorsInGroup(contract, collection, account string) (interface{}, error) {
	switch operators, st, err := hd.database.NFTOperators(contract, collection, account); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTOperatorsHal(contract, collection, account, *operators, st)
		if err != nil {
			return nil, err
		}
		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTOperatorsHal(contract, collection, account string, operators types.OperatorsBook, st base.State) (Hal, error) {
	h, err := hd.combineURL(HandlerPathNFTOperators, "contract", contract, "collection", collection, "account", account)
	if err != nil {
		return nil, err
	}

	hal := NewBaseHal(operators, NewHalLink(h, nil))

	return hal, nil
}
