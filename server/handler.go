package server

import (
	"crypto/md5"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
	"github.com/mr-tron/base58"
)

func (s *Server) addKeyHandler() http.HandlerFunc {
	req := struct {
		Key string `json:"key"`
	}{}

	resp := struct {
		Hash string `json:"hash"`
	}{}
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusBadRequest, err, "failed to decode request body")
			return
		}

		resp.Hash = hash(req.Key)
		// try get key form store
		if _, ok := s.store[resp.Hash]; ok {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusConflict, nil, "key already exists")
			return
		}
		s.store[resp.Hash] = key{
			PublicKey: req.Key,
			Secret:    "",
		}

		s.cleaner(resp.Hash)
		rest.RenderJSON(w, resp)
	}
}

func (s *Server) addSecretHandler() http.HandlerFunc {
	req := struct {
		Secret string `json:"secret"`
	}{}

	return func(w http.ResponseWriter, r *http.Request) {
		hashParam := chi.URLParam(r, "hash")
		if s.store[hashParam].Secret != "" {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusConflict, nil, "secret already exists")
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusBadRequest, err, "failed to decode request body")
			return
		}
		defer r.Body.Close()

		if _, ok := s.store[hashParam]; !ok {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusNotFound, nil, "key not found or expired")
			return
		}

		s.store[hashParam] = key{
			PublicKey: s.store[hashParam].PublicKey,
			Secret:    req.Secret,
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) getSecretHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashParam := chi.URLParam(r, "hash")

		if _, ok := s.store[hashParam]; !ok {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusNotFound, nil, "secret not found")
			return
		}
		rest.RenderJSON(w, s.store[hashParam])
	}
}

func (s *Server) deleteSecretHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashParam := chi.URLParam(r, "hash")

		if _, ok := s.store[hashParam]; !ok {
			rest.SendErrorJSON(w, r, log.Default(), http.StatusNotFound, nil, "secret not found")
			return
		}

		delete(s.store, hashParam)
		w.WriteHeader(http.StatusNoContent)
	}
}

// hash string md5 and base58 encode
func hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return base58.Encode(h.Sum(nil))
}
