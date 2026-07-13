package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("idempotency key conflicts with existing deposit")
)

type MemoryStore struct {
	mu          sync.Mutex
	nextID      int64
	users       map[string]domain.User
	addresses   map[string]domain.Address
	wallets     map[string]domain.Wallet
	deposits    map[string]domain.Deposit
	depositByTx map[string]string
	withdrawals map[string]domain.Withdrawal
	blacklist   map[string]domain.BlacklistEntry
	nonces      map[string]int64
	chainTxs    map[string]domain.ChainTx
	collections map[string]domain.CollectionTask
	auditLogs   []domain.AuditLog
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		users:       map[string]domain.User{},
		addresses:   map[string]domain.Address{},
		wallets:     map[string]domain.Wallet{},
		deposits:    map[string]domain.Deposit{},
		depositByTx: map[string]string{},
		withdrawals: map[string]domain.Withdrawal{},
		blacklist:   map[string]domain.BlacklistEntry{},
		nonces:      map[string]int64{},
		chainTxs:    map[string]domain.ChainTx{},
		collections: map[string]domain.CollectionTask{},
	}
	s.seedWallet(domain.ChainETH, "USDT", "0xhot-wallet-eth-usdt", domain.WalletHot, 1_000_000_000)
	s.seedWallet(domain.ChainETH, "USDT", "0xcold-wallet-eth-usdt", domain.WalletCold, 10_000_000_000)
	s.seedWallet(domain.ChainTRON, "USDT", "ThotWalletTronUSDT", domain.WalletHot, 1_000_000_000)
	return s
}

func (s *MemoryStore) seedWallet(chain domain.Chain, asset, address string, walletType domain.WalletType, balance int64) {
	id := s.newIDLocked("wallet")
	s.wallets[id] = domain.Wallet{ID: id, Chain: chain, Asset: asset, Address: address, WalletType: walletType, Balance: balance, CreatedAt: time.Now().UTC()}
}

func (s *MemoryStore) newIDLocked(prefix string) string {
	s.nextID++
	return fmt.Sprintf("%s_%06d", prefix, s.nextID)
}

func (s *MemoryStore) CreateUser(email string) domain.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	u := domain.User{ID: s.newIDLocked("user"), Email: email, CreatedAt: now}
	s.users[u.ID] = u
	s.addAuditLocked("system", "USER_CREATED", "user", u.ID, email)
	return u
}

func (s *MemoryStore) CreateAddress(userID string, chain domain.Chain) (domain.Address, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; !ok {
		return domain.Address{}, ErrNotFound
	}
	now := time.Now().UTC()
	addr := domain.Address{
		ID: s.newIDLocked("addr"), UserID: userID, Chain: chain,
		Address: fmt.Sprintf("%s-deposit-%s", chain, userID), WalletType: domain.WalletHot, CreatedAt: now,
	}
	s.addresses[addr.ID] = addr
	s.addAuditLocked("system", "ADDRESS_ALLOCATED", "address", addr.ID, addr.Address)
	return addr, nil
}

func (s *MemoryStore) AddressBelongsToUser(userID string, chain domain.Chain, address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.addresses {
		if a.UserID == userID && a.Chain == chain && a.Address == address {
			return true
		}
	}
	return false
}

func (s *MemoryStore) CreateDeposit(d domain.Deposit) (domain.Deposit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := string(d.Chain) + ":" + d.TxHash
	if existingID, ok := s.depositByTx[key]; ok {
		existing := s.deposits[existingID]
		if existing.UserID == d.UserID && existing.Asset == d.Asset && existing.Address == d.Address && existing.Amount == d.Amount {
			return existing, nil
		}
		return domain.Deposit{}, ErrConflict
	}
	now := time.Now().UTC()
	d.ID = s.newIDLocked("dep")
	d.Status = domain.DepositCredited
	d.Confirmations = 12
	d.CreatedAt = now
	d.UpdatedAt = now
	s.deposits[d.ID] = d
	s.depositByTx[key] = d.ID
	s.addAuditLocked("chain-indexer", "DEPOSIT_CREDITED", "deposit", d.ID, d.TxHash)
	return d, nil
}

func (s *MemoryStore) ListDeposits(userID string) []domain.Deposit {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []domain.Deposit{}
	for _, d := range s.deposits {
		if d.UserID == userID {
			out = append(out, d)
		}
	}
	return out
}

func (s *MemoryStore) CreateWithdrawal(w domain.Withdrawal) domain.Withdrawal {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	w.ID = s.newIDLocked("wd")
	w.RequestedAt = now
	w.UpdatedAt = now
	s.withdrawals[w.ID] = w
	s.addAuditLocked(w.UserID, "WITHDRAWAL_REQUESTED", "withdrawal", w.ID, w.ToAddress)
	return w
}

func (s *MemoryStore) GetWithdrawal(id string) (domain.Withdrawal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	w, ok := s.withdrawals[id]
	if !ok {
		return domain.Withdrawal{}, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) UpdateWithdrawal(w domain.Withdrawal, actor, action string) domain.Withdrawal {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.UpdatedAt = time.Now().UTC()
	s.withdrawals[w.ID] = w
	s.addAuditLocked(actor, action, "withdrawal", w.ID, string(w.Status))
	return w
}

func (s *MemoryStore) AddBlacklist(chain domain.Chain, address, reason string) domain.BlacklistEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := domain.BlacklistEntry{ID: s.newIDLocked("blk"), Chain: chain, Address: address, Reason: reason, CreatedAt: time.Now().UTC()}
	s.blacklist[s.blacklistKey(chain, address)] = entry
	s.addAuditLocked("admin", "BLACKLIST_ADDED", "blacklist", entry.ID, address)
	return entry
}

func (s *MemoryStore) IsBlacklisted(chain domain.Chain, address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.blacklist[s.blacklistKey(chain, address)]
	return ok
}

func (s *MemoryStore) NextNonce(chain domain.Chain, fromAddress string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := string(chain) + ":" + fromAddress
	s.nonces[key]++
	return s.nonces[key]
}

func (s *MemoryStore) HotWallet(chain domain.Chain, asset string) (domain.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, w := range s.wallets {
		if w.Chain == chain && w.Asset == asset && w.WalletType == domain.WalletHot {
			return w, nil
		}
	}
	return domain.Wallet{}, ErrNotFound
}

func (s *MemoryStore) ColdWallet(chain domain.Chain, asset string) (domain.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, w := range s.wallets {
		if w.Chain == chain && w.Asset == asset && w.WalletType == domain.WalletCold {
			return w, nil
		}
	}
	return domain.Wallet{}, ErrNotFound
}

func (s *MemoryStore) CreateChainTx(tx domain.ChainTx) domain.ChainTx {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	tx.ID = s.newIDLocked("tx")
	tx.TxHash = fmt.Sprintf("0xsimulated%s", tx.ID)
	tx.Status = domain.ChainTxBroadcast
	tx.CreatedAt = now
	tx.UpdatedAt = now
	s.chainTxs[tx.ID] = tx
	s.addAuditLocked("wallet-signer", "CHAIN_TX_BROADCASTED", "chain_tx", tx.ID, tx.TxHash)
	return tx
}

func (s *MemoryStore) ConfirmChainTx(id string) (domain.ChainTx, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, ok := s.chainTxs[id]
	if !ok {
		return domain.ChainTx{}, ErrNotFound
	}
	tx.Status = domain.ChainTxConfirmed
	tx.UpdatedAt = time.Now().UTC()
	s.chainTxs[id] = tx
	if tx.BusinessType == "withdrawal" {
		w := s.withdrawals[tx.BusinessID]
		w.Status = domain.WithdrawalConfirmed
		w.UpdatedAt = tx.UpdatedAt
		s.withdrawals[w.ID] = w
	}
	s.addAuditLocked("chain-indexer", "CHAIN_TX_CONFIRMED", "chain_tx", id, tx.TxHash)
	return tx, nil
}

func (s *MemoryStore) CreateCollection(chain domain.Chain, asset string, threshold int64) (domain.CollectionTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var from *domain.Wallet
	var to *domain.Wallet
	for _, w := range s.wallets {
		w := w
		if w.Chain == chain && w.Asset == asset && w.WalletType == domain.WalletHot && w.Balance >= threshold {
			from = &w
		}
		if w.Chain == chain && w.Asset == asset && w.WalletType == domain.WalletCold {
			to = &w
		}
	}
	if from == nil || to == nil {
		return domain.CollectionTask{}, ErrNotFound
	}
	task := domain.CollectionTask{
		ID: s.newIDLocked("col"), Chain: chain, Asset: asset, FromAddress: from.Address,
		ToAddress: to.Address, Amount: from.Balance - threshold/2, Status: domain.CollectionPlanned, CreatedAt: time.Now().UTC(),
	}
	s.collections[task.ID] = task
	s.addAuditLocked("treasury", "COLLECTION_PLANNED", "collection", task.ID, task.FromAddress+"->"+task.ToAddress)
	return task, nil
}

func (s *MemoryStore) AuditLogs() []domain.AuditLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.AuditLog, len(s.auditLogs))
	copy(out, s.auditLogs)
	return out
}

func (s *MemoryStore) Stats() (withdrawals, deposits, auditLogs int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.withdrawals), len(s.deposits), len(s.auditLogs)
}

func (s *MemoryStore) blacklistKey(chain domain.Chain, address string) string {
	return string(chain) + ":" + address
}

func (s *MemoryStore) addAuditLocked(actor, action, entity, entityID, detail string) {
	log := domain.AuditLog{
		ID: s.newIDLocked("audit"), Actor: actor, Action: action, Entity: entity,
		EntityID: entityID, Detail: detail, CreatedAt: time.Now().UTC(),
	}
	s.auditLogs = append(s.auditLogs, log)
}
