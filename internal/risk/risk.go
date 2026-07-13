package risk

import (
	"fmt"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/domain"
)

type Store interface {
	IsBlacklisted(chain domain.Chain, address string) bool
}

type Engine struct {
	store             Store
	manualReviewLimit int64
	rejectLimit       int64
}

type Decision struct {
	Allowed      bool   `json:"allowed"`
	NeedReview   bool   `json:"need_review"`
	RejectReason string `json:"reject_reason,omitempty"`
}

func NewEngine(store Store) *Engine {
	return &Engine{store: store, manualReviewLimit: 5_000_000, rejectLimit: 50_000_000}
}

func (e *Engine) CheckWithdrawal(chain domain.Chain, toAddress string, amount int64) Decision {
	if e.store.IsBlacklisted(chain, toAddress) {
		return Decision{Allowed: false, RejectReason: "destination address is blacklisted"}
	}
	if amount <= 0 {
		return Decision{Allowed: false, RejectReason: "amount must be positive"}
	}
	if amount > e.rejectLimit {
		return Decision{Allowed: false, RejectReason: fmt.Sprintf("amount exceeds hard limit %d", e.rejectLimit)}
	}
	if amount >= e.manualReviewLimit {
		return Decision{Allowed: true, NeedReview: true, RejectReason: "large withdrawal requires manual review"}
	}
	return Decision{Allowed: true, NeedReview: true}
}
