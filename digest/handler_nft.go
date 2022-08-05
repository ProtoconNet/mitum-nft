package digest

import (
	"net/http"
	"strings"
	"time"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/util"
)

func (hd *Handlers) handleNFT(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var id string
	s, found := mux.Vars(r)["id"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty id"), http.StatusBadRequest)

		return
	}
	id = s
	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTInGroup(id)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)
		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTInGroup(id string) (interface{}, error) {
	switch va, _, _, err := hd.database.NFT(id); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTHal(va)
		if err != nil {
			return nil, err
		}
		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTHal(va NFTValue) (Hal, error) {
	hinted := va.nft.ID().String()
	h, err := hd.combineURL(HandlerPathNFT, "id", hinted)
	if err != nil {
		return nil, err
	}

	hal := NewBaseHal(va, NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleNFTCollection(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var symbol string
	s, found := mux.Vars(r)["symbol"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty symbol"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty symbol"), http.StatusBadRequest)

		return
	}
	symbol = s
	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTCollectionInGroup(symbol)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)
		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTCollectionInGroup(symbol string) (interface{}, error) {
	switch va, _, _, err := hd.database.NFTCollection(symbol); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTCollectionHal(va)
		if err != nil {
			return nil, err
		}
		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTCollectionHal(va nft.Design) (Hal, error) {
	hinted := va.Symbol().String()
	h, err := hd.combineURL(HandlerPathNFTCollection, "symbol", hinted)
	if err != nil {
		return nil, err
	}

	hal := NewBaseHal(va, NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleCollectionNFTs(w http.ResponseWriter, r *http.Request) {
	var symbol string
	s, found := mux.Vars(r)["symbol"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty symbol"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty symbol"), http.StatusBadRequest)

		return
	}
	symbol = s
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := CacheKey(
		r.URL.Path, stringOffsetQuery(offset),
		stringBoolQuery("reverse", reverse),
	)

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleCollectionNFTsInGroup(symbol, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Error().Err(err).Str("symbol", symbol).Msg("failed to get nfts")
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

func (hd *Handlers) handleCollectionNFTsInGroup(
	symbol string,
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
		symbol, reverse, offset, limit,
		func(_ string, va NFTValue) (bool, error) {
			hal, err := hd.buildNFTHal(va)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, err
	} else if len(vas) < 1 {
		return nil, false, util.NotFoundError.Errorf("nfts not found")
	}

	i, err := hd.buildCollectionNFTsHal(symbol, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.enc.Marshal(i)
	return b, int64(len(vas)) == limit, err
}
