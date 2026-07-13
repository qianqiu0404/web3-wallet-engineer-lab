package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/httpapi"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/service"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/store"
)

func main() {
	addr := getenv("HTTP_ADDR", ":8080")
	svc := service.NewWalletService(store.NewMemoryStore())
	handler := httpapi.NewHandler(svc)

	log.Printf("web3 wallet engineer lab listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
