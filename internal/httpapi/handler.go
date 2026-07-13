package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/domain"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/service"
)

type Handler struct {
	svc *service.WalletService
	mux *http.ServeMux
}

func NewHandler(svc *service.WalletService) *Handler {
	h := &Handler{svc: svc, mux: http.NewServeMux()}
	h.routes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("GET /healthz", h.health)
	h.mux.HandleFunc("GET /metrics", h.metrics)
	h.mux.HandleFunc("POST /api/users", h.createUser)
	h.mux.HandleFunc("POST /api/users/{id}/addresses", h.createAddress)
	h.mux.HandleFunc("GET /api/users/{id}/deposits", h.listDeposits)
	h.mux.HandleFunc("POST /api/deposits/simulate", h.simulateDeposit)
	h.mux.HandleFunc("POST /api/withdrawals", h.requestWithdrawal)
	h.mux.HandleFunc("POST /api/admin/withdrawals/{id}/approve", h.approveWithdrawal)
	h.mux.HandleFunc("POST /api/admin/withdrawals/{id}/reject", h.rejectWithdrawal)
	h.mux.HandleFunc("POST /api/admin/blacklist", h.addBlacklist)
	h.mux.HandleFunc("POST /api/admin/collection-tasks", h.planCollection)
	h.mux.HandleFunc("GET /api/admin/audit-logs", h.auditLogs)
	h.mux.HandleFunc("POST /api/chain/tx/{id}/confirm", h.confirmTx)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) metrics(w http.ResponseWriter, _ *http.Request) {
	withdrawals, deposits, auditLogs := h.svc.Store().Stats()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP wallet_withdrawals_total Number of withdrawal records\n")
	fmt.Fprintf(w, "# TYPE wallet_withdrawals_total counter\nwallet_withdrawals_total %d\n", withdrawals)
	fmt.Fprintf(w, "# HELP wallet_deposits_total Number of deposit records\n")
	fmt.Fprintf(w, "# TYPE wallet_deposits_total counter\nwallet_deposits_total %d\n", deposits)
	fmt.Fprintf(w, "# HELP wallet_audit_logs_total Number of audit logs\n")
	fmt.Fprintf(w, "# TYPE wallet_audit_logs_total counter\nwallet_audit_logs_total %d\n", auditLogs)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	writeJSON(w, http.StatusCreated, h.svc.CreateUser(req.Email))
}

func (h *Handler) createAddress(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Chain domain.Chain `json:"chain"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	addr, err := h.svc.CreateAddress(r.PathValue("id"), req.Chain)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, addr)
}

func (h *Handler) simulateDeposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string       `json:"user_id"`
		Chain   domain.Chain `json:"chain"`
		Asset   string       `json:"asset"`
		Address string       `json:"address"`
		Amount  int64        `json:"amount"`
		TxHash  string       `json:"tx_hash"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	dep, err := h.svc.SimulateDeposit(req.UserID, req.Chain, req.Asset, req.Address, req.Amount, req.TxHash)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, dep)
}

func (h *Handler) listDeposits(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Store().ListDeposits(r.PathValue("id")))
}

func (h *Handler) requestWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string       `json:"user_id"`
		Chain     domain.Chain `json:"chain"`
		Asset     string       `json:"asset"`
		ToAddress string       `json:"to_address"`
		Amount    int64        `json:"amount"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	writeJSON(w, http.StatusCreated, h.svc.RequestWithdrawal(req.UserID, req.Chain, req.Asset, req.ToAddress, req.Amount))
}

func (h *Handler) approveWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Operator string `json:"operator"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	if req.Operator == "" {
		req.Operator = "admin"
	}
	result, err := h.svc.ApproveWithdrawal(r.PathValue("id"), req.Operator)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) rejectWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Operator string `json:"operator"`
		Reason   string `json:"reason"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	result, err := h.svc.RejectWithdrawal(r.PathValue("id"), req.Operator, req.Reason)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) addBlacklist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Chain   domain.Chain `json:"chain"`
		Address string       `json:"address"`
		Reason  string       `json:"reason"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	writeJSON(w, http.StatusCreated, h.svc.AddBlacklist(req.Chain, req.Address, req.Reason))
}

func (h *Handler) planCollection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Chain     domain.Chain `json:"chain"`
		Asset     string       `json:"asset"`
		Threshold int64        `json:"threshold"`
	}
	if decode(w, r, &req) != nil {
		return
	}
	task, err := h.svc.PlanCollection(req.Chain, req.Asset, req.Threshold)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) auditLogs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Store().AuditLogs())
}

func (h *Handler) confirmTx(w http.ResponseWriter, r *http.Request) {
	tx, err := h.svc.Store().ConfirmChainTx(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tx)
}

func decode(w http.ResponseWriter, r *http.Request, out any) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		msg = http.StatusText(status)
	}
	writeJSON(w, status, map[string]string{"error": msg})
}
