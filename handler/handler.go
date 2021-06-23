package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"test.com/cache-server/cache"
	"test.com/cache-server/store"
	"test.com/cache-server/utils"
)

// Handler handles URL routes
// it has cacher and storer instances
type Handler struct {
	*mux.Router
	cache cache.Cacher
	store store.Storer
}

// NewHandler creates instance of Handler
func NewHandler(options ...Option) (*Handler, error) {
	h := &Handler{}

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.writeResponse(w, http.StatusNotFound,
			store.ApiResponse{Error: http.StatusText(http.StatusNotFound), Message: "path not found"})
	})

	router.HandleFunc("/contact/{contact_id}", h.GetContact).Methods("GET")
	router.HandleFunc("/contact", h.UpsertContact).Methods("PUT", "POST")

	h.Router = router

	for _, opt := range options {
		opt(h)
	}

	if h.cache == nil {
		// TODO, get proper cache options via config/env or pass through parameter
		cache, err := cache.NewRedisCache()
		if err != nil {
			return nil, err
		}
		h.cache = cache
	}
	if h.store == nil {
		// TODO, get proper store options via config/env or pass through parameter
		store, err := store.NewAutoPilotStore()
		if err != nil {
			return nil, err
		}
		h.store = store
	}

	return h, nil
}

// GetContact gets contact from cache service first
// if not exists, then get through storer api
func (h *Handler) GetContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contactID := params["contact_id"]

	ctx, cancel := context.WithTimeout(context.Background(), h.cache.ReadTimeout())
	defer cancel()

	var contact Contact

	key := utils.GetKeyFromParams("contact", contactID)
	err := h.cache.Get(ctx, key, &contact)
	if err == nil {
		log.Info("hit cache") // TODO, replace by metrics counter
		h.writeResponse(w, http.StatusOK, contact)
		return
	}

	if !errors.Is(err, cache.ErrorKeyNotExists) {
		log.WithField("key", key).WithError(err).Error("error getting key from cache")
		h.writeResponse(w, http.StatusInternalServerError,
			store.ApiResponse{Error: http.StatusText(http.StatusInternalServerError), Message: "error getting key from cache"})
		return
	}

	ctx, cancelStorer := context.WithTimeout(ctx, h.store.Timeout())
	defer cancelStorer()
	// retrieve from store
	resp := h.store.Get(ctx, key, &contact)
	if resp.StatusCode != http.StatusOK { // retrieve from store error
		h.writeResponse(w, resp.StatusCode, resp)
		return
	}

	// set the key
	err = h.cache.Set(ctx, key, contact, h.cache.TTL())
	if err != nil {
		log.WithField("key", key).WithError(err).Error("error setting key")
		h.writeResponse(w, http.StatusInternalServerError,
			&store.ApiResponse{Error: http.StatusText(http.StatusInternalServerError), Message: "error setting key to cache"})
		return
	}

	h.writeResponse(w, http.StatusOK, contact)
}

// UpsertContact add contacts though storer api,
// and update the cache
func (h *Handler) UpsertContact(w http.ResponseWriter, r *http.Request) {
	// get contact
	var contactMap map[string]Contact
	err := json.NewDecoder(r.Body).Decode(&contactMap)
	if err != nil {
		h.writeResponse(w, http.StatusBadRequest,
			store.ApiResponse{Error: http.StatusText(http.StatusBadRequest), Message: "error decoding request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.store.Timeout())
	defer cancel()

	contact := contactMap["contact"]

	key := utils.GetKeyFromParams("contact", contact.ID)

	// Use distribute lock that make sure the data queried is the same being upserted
	// Since the strategy of the update is the later update wins, so it was not implemented

	var resp Contact
	rsUpsert := h.store.Upsert(ctx, key, contactMap, &resp)
	if rsUpsert.StatusCode != 200 {
		h.writeResponse(w, rsUpsert.StatusCode, rsUpsert)
		return
	}

	if resp.ID == "" {
		h.writeResponse(w, http.StatusInternalServerError,
			store.ApiResponse{Error: http.StatusText(http.StatusInternalServerError), Message: "got invalid contact ID " + resp.ID})
		return
	}

	// DO a query of the result, since there can be custom fields in the input,
	// which will be converted by the storer api
	ctx, cancelQuery := context.WithTimeout(ctx, h.store.Timeout())
	defer cancelQuery()

	key = utils.GetKeyFromParams("contact", resp.ID)
	// retrieve from store
	rsGet := h.store.Get(ctx, key, &contact)
	if rsGet.StatusCode != http.StatusOK {
		h.writeResponse(w, rsGet.StatusCode, rsGet)
		return
	}

	// set the key in store
	ctx, cancelCacher := context.WithTimeout(ctx, h.cache.WriteTimeout())
	defer cancelCacher()

	err = h.cache.Set(ctx, key, contact, h.cache.TTL())
	if err != nil {
		h.writeResponse(w, http.StatusInternalServerError,
			store.ApiResponse{Error: http.StatusText(http.StatusInternalServerError), Message: "error setting cache"})
	}

	// return the result of upsert response
	h.writeResponse(w, http.StatusOK, resp)

	return
}

func (h *Handler) writeResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshalling json payload"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		log.WithError(err).Error("error writing payload")
	}
}
