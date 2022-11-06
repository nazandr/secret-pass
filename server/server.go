package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
)

type Server struct {
	port  string
	store map[string]key

	lifeSpan time.Duration
}

type key struct {
	PublicKey string `json:"public_key"`
	Secret    string `json:"secret"`
}

type lifeSpan struct {
	Lifespan int64 `json:"lifespan"`
}

func NewServer(port string, lifeSpan time.Duration) *Server {
	return &Server{
		port:     port,
		store:    make(map[string]key),
		lifeSpan: lifeSpan,
	}
}

func (s *Server) Run() error {
	log.Printf("[INFO] start srever on %s", s.port)
	log.Printf("[INFO] life span of key is %s", s.lifeSpan)

	httpServer := &http.Server{
		Addr:              s.port,
		Handler:           s.router(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ErrorLog:          log.ToStdLogger(log.Default(), "WARN"),
	}

	return httpServer.ListenAndServe()
}

func (s *Server) router() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RealIP, rest.Recoverer(log.Default()))
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: log.ToStdLogger(log.Default(), "INFO"),
	}))
	router.Use(middleware.Throttle(1000))

	router.Route("/api", func(r chi.Router) {
		r.Get("/lifespan", func(w http.ResponseWriter, r *http.Request) {
			rest.RenderJSON(w, lifeSpan{s.lifeSpan.Milliseconds()})
		}) // get lifespan
		r.Post("/key", s.addKeyHandler())                   // add new key and hash it
		r.Post("/secret/{hash}", s.addSecretHandler())      // add secret to key
		r.Get("/secret/{hash}", s.getSecretHandler())       // get public key and encoded secret(if exist) by hash
		r.Delete("/secret/{hash}", s.deleteSecretHandler()) // delete secret by hash
	})

	fs, err := rest.NewFileServer("/static", "./assets/static")
	if err != nil {
		log.Fatalf("[ERROR] failed to create file server: %s", err)
	}
	router.Handle("/static/*", fs)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/index.html")
	})
	router.Get("/{hash}", func(w http.ResponseWriter, r *http.Request) {
		paramHash := chi.URLParam(r, "hash")
		if _, ok := s.store[paramHash]; !ok {
			http.ServeFile(w, r, "./assets/404.html")
			return
		}
		if s.store[paramHash].Secret != "" {
			http.ServeFile(w, r, "./assets/decrypt.html")
			return
		}
		http.ServeFile(w, r, "./assets/secret.html")
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/404.html")
	})

	return router
}

// remove key from store after some time
func (s *Server) cleaner(key string) {
	ticker := time.NewTicker(s.lifeSpan)
	go func() {
		for range ticker.C {
			delete(s.store, key)
		}
	}()
}
