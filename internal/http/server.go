package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"log"
	nethttp "net/http" //  чтобы не конфликтовать с именем пакета
	"path"
	"strings"
	"time"

	"github.com/hitokirisss/order-service/internal/cache"
	"github.com/hitokirisss/order-service/internal/storage"
)

//go:embed web/*
var webFS embed.FS

type Server struct {
	http  *nethttp.Server
	store *storage.Store
	cache *cache.Cache
}

func New(addr string, st *storage.Store, c *cache.Cache, preload map[string]json.RawMessage) *Server {
	for k, v := range preload {
		c.Set(k, v)
	}

	mux := nethttp.NewServeMux()

	mux.HandleFunc("/healthz", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/order/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		id := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/order/"))
		if id == "" || strings.Contains(id, "/") {
			nethttp.Error(w, "bad order id", nethttp.StatusBadRequest)
			return
		}

		if raw, ok := c.Get(id); ok {
			writeJSONBytes(w, raw, nethttp.StatusOK)
			return
		}

		raw, err := st.GetOrderRawJSON(r.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				nethttp.Error(w, "not found", nethttp.StatusNotFound)
			} else {
				nethttp.Error(w, "internal error", nethttp.StatusInternalServerError)
			}
			return
		}
		c.Set(id, raw)
		writeJSONBytes(w, raw, nethttp.StatusOK)
	})

	mux.HandleFunc("/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		file := "web/index.html"
		if p := strings.TrimPrefix(r.URL.Path, "/"); p != "" && p != "index.html" {
			nethttp.Redirect(w, r, "/", nethttp.StatusFound)
			return
		}
		b, err := webFS.ReadFile(path.Clean(file))
		if err != nil {
			nethttp.Error(w, "not found", nethttp.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write(b)
	})

	srv := &nethttp.Server{
		Addr:              addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{http: srv, store: st, cache: c}
}

func (s *Server) Start() {
	go func() {
		log.Printf("http listening on %s", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

func writeJSONBytes(w nethttp.ResponseWriter, b []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

func logRequests(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
