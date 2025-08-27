package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "service is running")
	})

	fmt.Println("server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("server stopped with error:", err)
	}
}
