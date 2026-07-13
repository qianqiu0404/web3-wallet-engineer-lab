package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/service"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/store"
)

func TestHTTPHappyPath(t *testing.T) {
	handler := NewHandler(service.NewWalletService(store.NewMemoryStore()))

	userResp := postJSON(t, handler, "/api/users", map[string]any{"email": "candidate@example.com"})
	userID := userResp["id"].(string)

	addrResp := postJSON(t, handler, "/api/users/"+userID+"/addresses", map[string]any{"chain": "ETH"})
	address := addrResp["address"].(string)

	depResp := postJSON(t, handler, "/api/deposits/simulate", map[string]any{
		"user_id": userID, "chain": "ETH", "asset": "USDT", "address": address, "amount": 1000000, "tx_hash": "0xabc",
	})
	if depResp["status"] != "CREDITED" {
		t.Fatalf("deposit status = %v", depResp["status"])
	}

	wdResp := postJSON(t, handler, "/api/withdrawals", map[string]any{
		"user_id": userID, "chain": "ETH", "asset": "USDT", "to_address": "0xreceiver", "amount": 1000000,
	})
	wdID := wdResp["id"].(string)

	approved := postJSON(t, handler, "/api/admin/withdrawals/"+wdID+"/approve", map[string]any{"operator": "admin"})
	if approved["status"] != "BROADCASTED" {
		t.Fatalf("withdrawal status = %v", approved["status"])
	}
}

func postJSON(t *testing.T, handler http.Handler, path string, body map[string]any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code < 200 || rec.Code >= 300 {
		t.Fatalf("%s returned %d: %s", path, rec.Code, rec.Body.String())
	}
	var out map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}
