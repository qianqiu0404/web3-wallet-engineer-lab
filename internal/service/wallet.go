package service

import (
	"errors"

	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/domain"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/risk"
	"github.com/qianqiu0404/web3-wallet-engineer-lab/internal/store"
)

type WalletService struct {
	store *store.MemoryStore
	risk  *risk.Engine
}

func NewWalletService(s *store.MemoryStore) *WalletService {
	return &WalletService{store: s, risk: risk.NewEngine(s)}
}

func (s *WalletService) Store() *store.MemoryStore {
	return s.store
}

func (s *WalletService) CreateUser(email string) domain.User {
	return s.store.CreateUser(email)
}

func (s *WalletService) CreateAddress(userID string, chain domain.Chain) (domain.Address, error) {
	return s.store.CreateAddress(userID, chain)
}

func (s *WalletService) SimulateDeposit(userID string, chain domain.Chain, asset, address string, amount int64, txHash string) (domain.Deposit, error) {
	if !s.store.AddressBelongsToUser(userID, chain, address) {
		return domain.Deposit{}, errors.New("deposit address does not belong to user")
	}
	return s.store.CreateDeposit(domain.Deposit{UserID: userID, Chain: chain, Asset: asset, Address: address, Amount: amount, TxHash: txHash})
}

func (s *WalletService) RequestWithdrawal(userID string, chain domain.Chain, asset, toAddress string, amount int64) domain.Withdrawal {
	decision := s.risk.CheckWithdrawal(chain, toAddress, amount)
	w := domain.Withdrawal{UserID: userID, Chain: chain, Asset: asset, ToAddress: toAddress, Amount: amount}
	if !decision.Allowed {
		w.Status = domain.WithdrawalRiskRejected
		w.RiskReason = decision.RejectReason
		return s.store.CreateWithdrawal(w)
	}
	w.Status = domain.WithdrawalPendingReview
	w.RiskReason = decision.RejectReason
	return s.store.CreateWithdrawal(w)
}

func (s *WalletService) ApproveWithdrawal(id, operator string) (domain.Withdrawal, error) {
	w, err := s.store.GetWithdrawal(id)
	if err != nil {
		return domain.Withdrawal{}, err
	}
	if w.Status != domain.WithdrawalPendingReview {
		return domain.Withdrawal{}, errors.New("withdrawal is not pending review")
	}
	hot, err := s.store.HotWallet(w.Chain, w.Asset)
	if err != nil {
		return domain.Withdrawal{}, err
	}
	nonce := s.store.NextNonce(w.Chain, hot.Address)
	tx := s.store.CreateChainTx(domain.ChainTx{
		Chain: w.Chain, Asset: w.Asset, FromAddress: hot.Address, ToAddress: w.ToAddress,
		Amount: w.Amount, Nonce: nonce, BusinessType: "withdrawal", BusinessID: w.ID,
	})
	w.Status = domain.WithdrawalBroadcasted
	w.Nonce = &nonce
	w.ChainTxID = tx.ID
	return s.store.UpdateWithdrawal(w, operator, "WITHDRAWAL_APPROVED"), nil
}

func (s *WalletService) RejectWithdrawal(id, operator, reason string) (domain.Withdrawal, error) {
	w, err := s.store.GetWithdrawal(id)
	if err != nil {
		return domain.Withdrawal{}, err
	}
	if w.Status != domain.WithdrawalPendingReview {
		return domain.Withdrawal{}, errors.New("withdrawal is not pending review")
	}
	w.Status = domain.WithdrawalRejected
	w.RiskReason = reason
	return s.store.UpdateWithdrawal(w, operator, "WITHDRAWAL_REJECTED"), nil
}

func (s *WalletService) AddBlacklist(chain domain.Chain, address, reason string) domain.BlacklistEntry {
	return s.store.AddBlacklist(chain, address, reason)
}

func (s *WalletService) PlanCollection(chain domain.Chain, asset string, threshold int64) (domain.CollectionTask, error) {
	return s.store.CreateCollection(chain, asset, threshold)
}
