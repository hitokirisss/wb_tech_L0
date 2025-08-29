package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hitokirisss/order-service/internal/cache"
	"github.com/hitokirisss/order-service/internal/config"
	"github.com/hitokirisss/order-service/internal/http"
	"github.com/hitokirisss/order-service/internal/storage"
)

func main() {
	cfg := config.New()

	ctx := context.Background()
	st, err := storage.New(ctx, cfg.PGURL())
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer st.Close()

	c := cache.New()

	// предзагрузка кеша
	pre := map[string]json.RawMessage{}
	if cfg.CachePreload > 0 {
		m, err := st.LoadRecentRaw(ctx, cfg.CachePreload)
		if err != nil {
			log.Printf("preload cache: %v", err)
		} else {
			for k, v := range m {
				pre[k] = v
			}
			log.Printf("cache preloaded: %d", len(m))
		}
	}

	srv := http.New(cfg.HTTPAddr, st, c, pre)
	srv.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxSh, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxSh)
	log.Println("stopped")
}
