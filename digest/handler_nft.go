package digest

import (
	"net/http"
	"strings"
	"time"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

func (hd *Handlers) buildNFTHal(va nft.NFT) (Hal, error) {
	hinted := va.ID().String()
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
