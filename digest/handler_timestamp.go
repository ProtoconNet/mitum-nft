package digest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ProtoconNet/mitum-nft/v2/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleTimeStamp(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var contract string
	s, found := mux.Vars(r)["contract"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusBadRequest)

		return
	}
	contract = s

	var service string
	s, found = mux.Vars(r)["service"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusBadRequest)

		return
	}
	service = s

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleTimeStampInGroup(contract, service)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)

		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampInGroup(contract, service string) ([]byte, error) {
	var de timestamp.Design
	var st base.State

	de, st, err := hd.database.timestamp(contract, service)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStamp(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.enc.Marshal(i)
}

func (hd *Handlers) buildTimeStamp(contract string, de timestamp.Design, st base.State) (Hal, error) {
	h, err := hd.combineURL(HandlerPathTimeStampService, "contract", contract, "service", de.Service().String())
	if err != nil {
		return nil, err
	}

	var hal Hal
	hal = NewBaseHal(de, NewHalLink(h, nil))

	h, err = hd.combineURL(HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleTimeStampItem(w http.ResponseWriter, r *http.Request) {
	cachekey := CacheKeyPath(r)
	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var contract string
	s, found := mux.Vars(r)["contract"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusBadRequest)

		return
	}
	contract = s

	var service string
	s, found = mux.Vars(r)["service"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusBadRequest)

		return
	}
	service = s

	var project string
	s, found = mux.Vars(r)["project"]
	if !found {
		HTTP2ProblemWithError(w, errors.Errorf("empty project id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		HTTP2ProblemWithError(w, errors.Errorf("empty project id"), http.StatusBadRequest)

		return
	}
	project = s

	s, found = mux.Vars(r)["tid"]
	idx, err := parseIdxFromPath(s)
	if err != nil {
		HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleTimeStampItemInGroup(contract, service, project, idx)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)

		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampItemInGroup(contract, service, project string, idx uint64) ([]byte, error) {
	var it timestamp.TimeStampItem
	var st base.State

	it, st, err := hd.database.timestampItem(contract, service, project, idx)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStampItem(contract, service, it, st)
	if err != nil {
		return nil, err
	}
	return hd.enc.Marshal(i)
}

func (hd *Handlers) buildTimeStampItem(contract, service string, it timestamp.TimeStampItem, st base.State) (Hal, error) {
	h, err := hd.combineURL(HandlerPathTimeStampItem, "contract", contract, "service", service, "project", it.ProjectID(), "tid", strconv.FormatUint(it.TimestampID(), 10))
	if err != nil {
		return nil, err
	}

	var hal Hal
	hal = NewBaseHal(it, NewHalLink(h, nil))

	h, err = hd.combineURL(HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", NewHalLink(h, nil))
	}

	return hal, nil
}
