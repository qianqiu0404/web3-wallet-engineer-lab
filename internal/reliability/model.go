package reliability

import (
	"errors"
	"fmt"
)

var ErrIdempotencyConflict = errors.New("idempotency key reused with different payload")

type Withdrawal struct {
	RequestID       string
	PayloadHash     string
	Status          string
	FreezeCount     int
	CandidateHashes []string
	CanonicalHash   string
}

type WithdrawalBook struct {
	items map[string]*Withdrawal
}

func NewWithdrawalBook() *WithdrawalBook {
	return &WithdrawalBook{items: map[string]*Withdrawal{}}
}

func (b *WithdrawalBook) Submit(requestID, payloadHash string) (*Withdrawal, error) {
	if existing, ok := b.items[requestID]; ok {
		if existing.PayloadHash != payloadHash {
			return nil, ErrIdempotencyConflict
		}
		return existing, nil
	}
	w := &Withdrawal{RequestID: requestID, PayloadHash: payloadHash, Status: "review_pending", FreezeCount: 1}
	b.items[requestID] = w
	return w, nil
}

func (w *Withdrawal) MarkBroadcastUnknown(rawFingerprint string) error {
	if rawFingerprint == "" {
		return errors.New("raw transaction fingerprint is required")
	}
	w.Status = "broadcast_unknown"
	return nil
}

func (w *Withdrawal) TrackCandidate(txHash string) {
	for _, existing := range w.CandidateHashes {
		if existing == txHash {
			return
		}
	}
	w.CandidateHashes = append(w.CandidateHashes, txHash)
}

func (w *Withdrawal) ConfirmCanonical(txHash string) error {
	found := false
	for _, candidate := range w.CandidateHashes {
		if candidate == txHash {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("canonical hash %s is not tracked", txHash)
	}
	if w.CanonicalHash != "" && w.CanonicalHash != txHash {
		return errors.New("withdrawal already has a different canonical result")
	}
	w.CanonicalHash = txHash
	w.Status = "wallet_done"
	return nil
}

type Settlement struct {
	LedgerEntries int
	Notifications int
	Status        string
}

func (s *Settlement) ApplyCanonical() {
	if s.LedgerEntries == 0 {
		s.LedgerEntries = 1
	}
	s.Status = "wallet_done"
}

func (s *Settlement) Notify() {
	if s.Status == "wallet_done" && s.Notifications == 0 {
		s.Notifications = 1
	}
}

type DepositLedger struct {
	PendingAmount   int64
	AvailableAmount int64
	ReversalEntries int
}

func (d *DepositLedger) Observe(amount int64) { d.PendingAmount += amount }

func (d *DepositLedger) Reorg(amount int64) {
	if d.PendingAmount >= amount {
		d.PendingAmount -= amount
		d.ReversalEntries++
	}
}

type AuthorizationGate struct {
	RiskAvailable   bool
	SignerAvailable bool
}

func (g AuthorizationGate) CanSign() bool { return g.RiskAvailable && g.SignerAvailable }
