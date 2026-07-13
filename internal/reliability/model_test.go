package reliability

import (
	"errors"
	"testing"
)

func TestBaselineHasOneCanonicalOutcome(t *testing.T) {
	book := NewWithdrawalBook()
	w, err := book.Submit("wd-1", "payload-a")
	if err != nil {
		t.Fatal(err)
	}
	w.TrackCandidate("0xtx")
	if err := w.ConfirmCanonical("0xtx"); err != nil {
		t.Fatal(err)
	}
	if w.FreezeCount != 1 || w.CanonicalHash != "0xtx" || w.Status != "wallet_done" {
		t.Fatalf("unexpected baseline: %+v", w)
	}
}

func TestDuplicateWithdrawalIsIdempotentAndConflictsAreRejected(t *testing.T) {
	book := NewWithdrawalBook()
	first, err := book.Submit("wd-1", "payload-a")
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := book.Submit("wd-1", "payload-a")
	if err != nil {
		t.Fatal(err)
	}
	if duplicate != first || first.FreezeCount != 1 {
		t.Fatal("duplicate request created a second fund action")
	}
	if _, err := book.Submit("wd-1", "payload-b"); !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("conflict error = %v", err)
	}
}

func TestBroadcastUnknownPausesNewTransaction(t *testing.T) {
	w := &Withdrawal{Status: "signed"}
	if err := w.MarkBroadcastUnknown("sha256:raw"); err != nil {
		t.Fatal(err)
	}
	if w.Status != "broadcast_unknown" || len(w.CandidateHashes) != 0 {
		t.Fatalf("unknown result must pause: %+v", w)
	}
}

func TestCanonicalCompensationIsIdempotent(t *testing.T) {
	settlement := &Settlement{}
	settlement.ApplyCanonical()
	settlement.ApplyCanonical()
	settlement.Notify()
	settlement.Notify()
	if settlement.LedgerEntries != 1 || settlement.Notifications != 1 {
		t.Fatalf("duplicate side effect: %+v", settlement)
	}
}

func TestDepositReorgUsesReversalEntry(t *testing.T) {
	ledger := &DepositLedger{}
	ledger.Observe(100)
	ledger.Reorg(100)
	if ledger.PendingAmount != 0 || ledger.AvailableAmount != 0 || ledger.ReversalEntries != 1 {
		t.Fatalf("unexpected reorg result: %+v", ledger)
	}
}

func TestNonceReplacementHasOneCanonicalOutcome(t *testing.T) {
	w := &Withdrawal{}
	w.TrackCandidate("0xold")
	w.TrackCandidate("0xreplacement")
	if err := w.ConfirmCanonical("0xreplacement"); err != nil {
		t.Fatal(err)
	}
	if err := w.ConfirmCanonical("0xold"); err == nil {
		t.Fatal("second canonical result was accepted")
	}
}

func TestRiskAndSignerOutageFailsClosed(t *testing.T) {
	for _, gate := range []AuthorizationGate{{RiskAvailable: false, SignerAvailable: true}, {RiskAvailable: true, SignerAvailable: false}, {RiskAvailable: false, SignerAvailable: false}} {
		if gate.CanSign() {
			t.Fatalf("unavailable dependency must fail closed: %+v", gate)
		}
	}
	if !(AuthorizationGate{RiskAvailable: true, SignerAvailable: true}).CanSign() {
		t.Fatal("healthy dependencies should allow signing")
	}
}
