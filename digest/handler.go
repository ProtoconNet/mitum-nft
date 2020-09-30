package digest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/spikeekips/mitum/network"
	"github.com/spikeekips/mitum/util/encoder"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/logging"
	"golang.org/x/xerrors"
)

var (
	HTTP2EncoderHintHeader = http.CanonicalHeaderKey("x-mitum-encoder-hint")
	HALMimetype            = "application/hal+json; charset=utf-8"
)

var (
	HandlerPathNodeInfo          = `/`
	HandlerPathBlockByHeight     = `/block/{height:[0-9]+}`
	HandlerPathBlockByHash       = `/block/{hash:(?i)[0-9a-z][0-9a-z]+}`
	HandlerPathManifestByHeight  = `/block/{height:[0-9]+}/manifest`
	HandlerPathManifestByHash    = `/block/{hash:(?i)[0-9a-z][0-9a-z]+}/manifest`
	HandlerPathAccount           = `/account/{address:(?i)[0-9a-z][0-9a-z\-]+\-[a-z0-9]{4}\:[a-z0-9\.]*}`
	HandlerPathAccountOperations = `/account/{address:(?i)[0-9a-z][0-9a-z\-]+\-[a-z0-9]{4}\:[a-z0-9\.]*}/operations`
	HandlerPathOperation         = `/operation/{hash:(?i)[0-9a-z][0-9a-z]+}`
)

var (
	UnknownProblem     = NewProblem(DefaultProblemType, "unknown problem occurred")
	unknownProblemJSON []byte
)

var GlobalLimitList int64 = 10

func init() {
	if b, err := jsonenc.Marshal(UnknownProblem); err != nil {
		panic(err)
	} else {
		unknownProblemJSON = b
	}
}

type Handlers struct {
	*logging.Logging
	enc             encoder.Encoder
	storage         *Storage
	cache           Cache
	nodeInfoHandler network.NodeInfoHandler
	router          *mux.Router
	routes          map[ /* path */ string]*mux.Route
	limiter         func(string /* request type */) int64
}

func NewHandlers(enc encoder.Encoder, st *Storage, cache Cache) *Handlers {
	return &Handlers{
		Logging: logging.NewLogging(func(c logging.Context) logging.Emitter {
			return c.Str("module", "http2-handlers")
		}),
		enc:     enc,
		storage: st,
		cache:   cache,
		router:  mux.NewRouter(),
		routes:  map[string]*mux.Route{},
		limiter: defaultLimiter,
	}
}

func (hd *Handlers) Initialize() error {
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	hd.router.Use(cors)

	hd.router.Use(func() mux.MiddlewareFunc {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cr := NewCacheResponseWriter(hd.cache, w, r)

				next.ServeHTTP(cr, r)

				if err := cr.Cache(); err != nil {
					if xerrors.Is(err, SkipCacheError) {
						hd.Log().Verbose().Msg("cache skipped")
					} else {
						hd.Log().Error().Err(err).Msg("failed to cache")
					}
				} else {
					hd.Log().Verbose().Msg("cached")
				}
			})
		}
	}())

	hd.setHandlers()

	return nil
}

func (hd *Handlers) SetLimiter(f func(string) int64) *Handlers {
	hd.limiter = f

	return hd
}

func (hd *Handlers) Handler() http.Handler {
	return network.HTTPLogHandler(hd.router, hd.Log())
}

func (hd *Handlers) setHandlers() {
	_ = hd.setHandler(HandlerPathManifestByHeight, hd.handleManifestByHeight).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathManifestByHash, hd.handleManifestByHash).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathBlockByHeight, hd.handleBlock).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathBlockByHash, hd.handleBlock).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathAccount, hd.handleAccount).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathAccountOperations, hd.handleAccountOperations).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathOperation, hd.handleOperation).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathNodeInfo, hd.handleNodeInfo).
		Methods(http.MethodOptions, "GET")
}

func (hd *Handlers) setHandler(prefix string, handler network.HTTPHandlerFunc) *mux.Route {
	var name string
	if prefix == "" || prefix == "/" {
		name = "root"
	} else {
		name = prefix
	}

	var route *mux.Route
	if r := hd.router.Get(name); r != nil {
		route = r
	} else {
		route = hd.router.Name(name)
	}

	route = route.
		Path(prefix).
		HandlerFunc(handler)

	hd.routes[prefix] = route

	return route
}

func (hd *Handlers) stream(w http.ResponseWriter, bufsize int, status int) (*jsoniter.Stream, func()) {
	w.Header().Set(HTTP2EncoderHintHeader, hd.enc.Hint().String())
	w.Header().Set("Content-Type", HALMimetype)

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	stream := jsoniter.NewStream(HALJSONConfigDefault, w, bufsize)
	return stream, func() {
		if err := stream.Flush(); err != nil {
			hd.Log().Error().Err(err).Msg("failed to straem thru jsoniterator")

			hd.problemWithError(w, err, http.StatusInternalServerError)
		}
	}
}

func (hd *Handlers) combineURL(path string, pairs ...string) (string, error) {
	if len(pairs)%2 != 0 {
		return "", xerrors.Errorf("failed to combine url; uneven pairs to combine url")
	}

	if u, err := hd.routes[path].URLPath(pairs...); err != nil {
		return "", xerrors.Errorf("failed to combine url: %w", err)
	} else {
		return u.String(), nil
	}
}

func (hd *Handlers) notSupported(w http.ResponseWriter) {
	hd.writePoblem(w, NewProblemFromError(xerrors.Errorf("not supported")), http.StatusInternalServerError)
}

func (hd *Handlers) problemWithError(w http.ResponseWriter, err error, status int) {
	hd.writePoblem(w, NewProblemFromError(err), status)
}

func (hd *Handlers) writePoblem(w http.ResponseWriter, pr Problem, status int) {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.Header().Set(HTTP2EncoderHintHeader, hd.enc.Hint().String())
	w.Header().Set("Content-Type", ProblemMimetype)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	var output []byte
	if b, err := jsonenc.Marshal(pr); err != nil {
		hd.Log().Error().Err(err).Interface("problem", pr).Msg("failed to marshal problem, UnknownProblem will be used")

		output = unknownProblemJSON
	} else {
		output = b
	}

	w.WriteHeader(status)
	_, _ = w.Write(output)
}

func (hd *Handlers) writeHal(w http.ResponseWriter, hal Hal, status int) { // nolint:unparam
	stream, flush := hd.stream(w, 1, status)
	defer flush()

	stream.WriteVal(hal)
}

func (hd *Handlers) writeCache(w http.ResponseWriter, key string, expire time.Duration) {
	if cw, ok := w.(*CacheResponseWriter); ok {
		_ = cw.SetKey(key).SetExpire(expire)
	}
}

func cacheKeyPath(r *http.Request) string {
	return r.URL.Path
}

func cacheKey(key string, s ...string) string {
	return fmt.Sprintf("%s-%s", key, strings.Join(s, ","))
}

func defaultLimiter(string) int64 {
	return GlobalLimitList
}