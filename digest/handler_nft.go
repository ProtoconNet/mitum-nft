package digest

import (
	"net/http"
	"strconv"
	"time"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var (
	HandlerPathNFTAllApproved = `/nft/{contract:(?i)` + ctypes.REStringAddressString + `}/account/{address:(?i)` + ctypes.REStringAddressString + `}/allapproved` // revive:disable-line:line-length-limit
	HandlerPathNFTCollection  = `/nft/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathNFT            = `/nft/{contract:(?i)` + ctypes.REStringAddressString + `}/nftidx/{nft_idx:[0-9]+}`
	HandlerPathNFTs           = `/nft/{contract:(?i)` + ctypes.REStringAddressString + `}/nfts`
)

func SetHandlers(hd *cdigest.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathNFTCollection, HandleNFTCollection, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathNFTs, HandleNFTs, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathNFTAllApproved, HandleNFTOperators, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathNFT, HandleNFT, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandleNFT(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	id, err, status := cdigest.ParseRequest(w, r, "nft_idx")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleNFTInGroup(hd, contract, id)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleNFTInGroup(hd *cdigest.Handlers, contract, id string) (interface{}, error) {
	switch nft, err := NFT(hd.Database(), contract, id); {
	case err != nil:
		return nil, err
	default:
		hal, err := buildNFTHal(hd, contract, *nft)
		if err != nil {
			return nil, err
		}
		return hd.Encoder().Marshal(hal)
	}
}

func buildNFTHal(hd *cdigest.Handlers, contract string, nft types.NFT) (cdigest.Hal, error) {
	h, err := hd.CombineURL(HandlerPathNFT, "contract", contract, "nft_idx", strconv.FormatUint(nft.ID(), 10))
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(nft, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func HandleNFTCollection(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleNFTCollectionInGroup(hd, contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleNFTCollectionInGroup(hd *cdigest.Handlers, contract string) (interface{}, error) {
	switch design, err := NFTCollection(hd.Database(), contract); {
	case err != nil:
		return nil, err
	default:
		hal, err := buildNFTCollectionHal(hd, contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.Encoder().Marshal(hal)
	}
}

func buildNFTCollectionHal(hd *cdigest.Handlers, contract string, design types.Design) (cdigest.Hal, error) {
	h, err := hd.CombineURL(HandlerPathNFTCollection, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(design, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func HandleNFTs(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	limit := cdigest.ParseLimitQuery(r.URL.Query().Get("limit"))
	offset := cdigest.ParseStringQuery(r.URL.Query().Get("offset"))
	reverse := cdigest.ParseBoolQuery(r.URL.Query().Get("reverse"))
	facthash := cdigest.ParseStringQuery(r.URL.Query().Get("facthash"))

	cachekey := cdigest.CacheKey(
		r.URL.Path, cdigest.StringOffsetQuery(offset),
		cdigest.StringBoolQuery("reverse", reverse),
	)

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		i, filled, err := handleNFTsInGroup(hd, contract, facthash, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("contract", contract).Msg("failed to get nfts")
		cdigest.HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, b, http.StatusOK)

	if !shared {
		expire := hd.ExpireNotFilled()
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		cdigest.HTTP2WriteCache(w, cachekey, expire)
	}
}

func handleNFTsInGroup(
	hd *cdigest.Handlers,
	contract, facthash, offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.ItemsLimiter("collection-nfts")
	} else {
		limit = l
	}

	var vas []cdigest.Hal
	if err := NFTsByCollection(
		hd.Database(), contract, facthash, offset, reverse, limit,
		func(nft types.NFT, st base.State) (bool, error) {
			hal, err := buildNFTHal(hd, contract, nft)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, util.ErrNotFound.WithMessage(err, "nft tokens by contract, %s", contract)
	} else if len(vas) < 1 {
		return nil, false, util.ErrNotFound.Errorf("nft tokens by contract, %s", contract)
	}

	i, err := buildNFTsHal(hd, contract, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.Encoder().Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func buildNFTsHal(
	hd *cdigest.Handlers,
	contract string,
	vas []cdigest.Hal,
	offset string,
	reverse bool,
) (cdigest.Hal, error) {
	baseSelf, err := hd.CombineURL(HandlerPathNFTs, "contract", contract)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringOffsetQuery(offset))
	}
	if reverse {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", reverse))
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(vas, cdigest.NewHalLink(self, nil))

	h, err := hd.CombineURL(HandlerPathNFTCollection, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("collection", cdigest.NewHalLink(h, nil))

	var nextoffset string

	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(types.NFT)
		nextoffset = strconv.FormatUint(va.ID(), 10)
	}

	if len(nextoffset) > 0 {
		next := baseSelf
		next = cdigest.AddQueryValue(next, cdigest.StringOffsetQuery(nextoffset))

		if reverse {
			next = cdigest.AddQueryValue(next, cdigest.StringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", cdigest.NewHalLink(next, nil))
	}

	hal = hal.AddLink(
		"reverse",
		cdigest.NewHalLink(
			cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", !reverse)),
			nil,
		),
	)

	return hal, nil
}

func HandleNFTOperators(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := cdigest.ParseRequest(w, r, "address")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleNFTOperatorsInGroup(hd, contract, account)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleNFTOperatorsInGroup(hd *cdigest.Handlers, contract, account string) (interface{}, error) {
	switch operators, err := NFTOperators(hd.Database(), contract, account); {
	case err != nil:
		return nil, err
	default:
		hal, err := buildNFTOperatorsHal(hd, contract, account, *operators)
		if err != nil {
			return nil, err
		}
		return hd.Encoder().Marshal(hal)
	}
}

func buildNFTOperatorsHal(hd *cdigest.Handlers, contract, account string, operators types.AllApprovedBook) (cdigest.Hal, error) {
	h, err := hd.CombineURL(HandlerPathNFTAllApproved, "contract", contract, "address", account)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(operators, cdigest.NewHalLink(h, nil))

	return hal, nil
}
