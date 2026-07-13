package service

import (
	"errors"
	"testing"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/domain"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/store"
)

func TestWithdrawalApprovalAllocatesNonceAndChainTx(t *testing.T) {
	svc := NewWalletService(store.NewMemoryStore())

	user := svc.CreateUser("candidate@example.com")
	w := svc.RequestWithdrawal(user.ID, domain.ChainETH, "USDT", "0xreceiver", 1_000_000)
	if w.Status != domain.WithdrawalPendingReview {
		t.Fatalf("status = %s, want %s", w.Status, domain.WithdrawalPendingReview)
	}

	approved, err := svc.ApproveWithdrawal(w.ID, "risk-admin")
	if err != nil {
		t.Fatalf("approve withdrawal: %v", err)
	}
	if approved.Status != domain.WithdrawalBroadcasted {
		t.Fatalf("status = %s, want %s", approved.Status, domain.WithdrawalBroadcasted)
	}
	if approved.Nonce == nil || *approved.Nonce != 1 {
		t.Fatalf("nonce = %v, want 1", approved.Nonce)
	}
	if approved.ChainTxID == "" {
		t.Fatal("chain tx id should be set")
	}
}

func TestBlacklistedWithdrawalIsRiskRejected(t *testing.T) {
	svc := NewWalletService(store.NewMemoryStore())
	user := svc.CreateUser("candidate@example.com")

	svc.AddBlacklist(domain.ChainETH, "0xbad", "sanctions screening hit")
	w := svc.RequestWithdrawal(user.ID, domain.ChainETH, "USDT", "0xbad", 1_000_000)

	if w.Status != domain.WithdrawalRiskRejected {
		t.Fatalf("status = %s, want %s", w.Status, domain.WithdrawalRiskRejected)
	}
	if w.RiskReason == "" {
		t.Fatal("risk reason should be present")
	}
}

func TestDepositRequiresAllocatedAddress(t *testing.T) {
	svc := NewWalletService(store.NewMemoryStore())
	user := svc.CreateUser("candidate@example.com")

	if _, err := svc.SimulateDeposit(user.ID, domain.ChainETH, "USDT", "0xunknown", 1_000_000, "0xtx"); err == nil {
		t.Fatal("expected deposit into unknown address to fail")
	}

	addr, err := svc.CreateAddress(user.ID, domain.ChainETH)
	if err != nil {
		t.Fatalf("create address: %v", err)
	}
	dep, err := svc.SimulateDeposit(user.ID, domain.ChainETH, "USDT", addr.Address, 1_000_000, "0xtx")
	if err != nil {
		t.Fatalf("simulate deposit: %v", err)
	}
	if dep.Status != domain.DepositCredited {
		t.Fatalf("status = %s, want %s", dep.Status, domain.DepositCredited)
	}
}

func TestDuplicateDepositIsIdempotentAndConflictingPayloadIsRejected(t *testing.T) {
	svc := NewWalletService(store.NewMemoryStore())
	user := svc.CreateUser("candidate@example.com")
	addr, err := svc.CreateAddress(user.ID, domain.ChainETH)
	if err != nil {
		t.Fatalf("create address: %v", err)
	}

	first, err := svc.SimulateDeposit(user.ID, domain.ChainETH, "USDT", addr.Address, 1_000_000, "0xsame")
	if err != nil {
		t.Fatalf("first deposit: %v", err)
	}
	duplicate, err := svc.SimulateDeposit(user.ID, domain.ChainETH, "USDT", addr.Address, 1_000_000, "0xsame")
	if err != nil {
		t.Fatalf("duplicate deposit: %v", err)
	}
	if duplicate.ID != first.ID {
		t.Fatalf("duplicate id = %s, want original %s", duplicate.ID, first.ID)
	}

	if _, err := svc.SimulateDeposit(user.ID, domain.ChainETH, "USDT", addr.Address, 2_000_000, "0xsame"); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("conflicting deposit error = %v, want ErrConflict", err)
	}
	if got := len(svc.Store().ListDeposits(user.ID)); got != 1 {
		t.Fatalf("deposit count = %d, want 1", got)
	}
}
