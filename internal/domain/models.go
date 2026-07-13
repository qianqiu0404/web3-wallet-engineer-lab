package domain

import "time"

type Chain string

const (
	ChainETH  Chain = "ETH"
	ChainBTC  Chain = "BTC"
	ChainTRON Chain = "TRON"
)

type WalletType string

const (
	WalletHot  WalletType = "HOT"
	WalletCold WalletType = "COLD"
)

type WithdrawalStatus string

const (
	WithdrawalPendingReview WithdrawalStatus = "PENDING_REVIEW"
	WithdrawalRiskRejected  WithdrawalStatus = "RISK_REJECTED"
	WithdrawalRejected      WithdrawalStatus = "REJECTED"
	WithdrawalApproved      WithdrawalStatus = "APPROVED"
	WithdrawalBroadcasted   WithdrawalStatus = "BROADCASTED"
	WithdrawalConfirmed     WithdrawalStatus = "CONFIRMED"
	WithdrawalFailed        WithdrawalStatus = "FAILED"
)

type DepositStatus string

const (
	DepositSeen      DepositStatus = "SEEN"
	DepositConfirmed DepositStatus = "CONFIRMED"
	DepositCredited  DepositStatus = "CREDITED"
)

type ChainTxStatus string

const (
	ChainTxCreated   ChainTxStatus = "CREATED"
	ChainTxBroadcast ChainTxStatus = "BROADCASTED"
	ChainTxConfirmed ChainTxStatus = "CONFIRMED"
	ChainTxFailed    ChainTxStatus = "FAILED"
)

type CollectionStatus string

const (
	CollectionPlanned   CollectionStatus = "PLANNED"
	CollectionBroadcast CollectionStatus = "BROADCASTED"
	CollectionDone      CollectionStatus = "DONE"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Address struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Chain      Chain      `json:"chain"`
	Address    string     `json:"address"`
	WalletType WalletType `json:"wallet_type"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Wallet struct {
	ID         string     `json:"id"`
	Chain      Chain      `json:"chain"`
	Asset      string     `json:"asset"`
	Address    string     `json:"address"`
	WalletType WalletType `json:"wallet_type"`
	Balance    int64      `json:"balance"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Deposit struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	Chain         Chain         `json:"chain"`
	Asset         string        `json:"asset"`
	Address       string        `json:"address"`
	Amount        int64         `json:"amount"`
	TxHash        string        `json:"tx_hash"`
	Confirmations int           `json:"confirmations"`
	Status        DepositStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Withdrawal struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	Chain       Chain            `json:"chain"`
	Asset       string           `json:"asset"`
	ToAddress   string           `json:"to_address"`
	Amount      int64            `json:"amount"`
	Status      WithdrawalStatus `json:"status"`
	RiskReason  string           `json:"risk_reason,omitempty"`
	Nonce       *int64           `json:"nonce,omitempty"`
	ChainTxID   string           `json:"chain_tx_id,omitempty"`
	RequestedAt time.Time        `json:"requested_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type BlacklistEntry struct {
	ID        string    `json:"id"`
	Chain     Chain     `json:"chain"`
	Address   string    `json:"address"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type ChainTx struct {
	ID           string        `json:"id"`
	Chain        Chain         `json:"chain"`
	Asset        string        `json:"asset"`
	FromAddress  string        `json:"from_address"`
	ToAddress    string        `json:"to_address"`
	Amount       int64         `json:"amount"`
	Nonce        int64         `json:"nonce"`
	TxHash       string        `json:"tx_hash"`
	Status       ChainTxStatus `json:"status"`
	BusinessType string        `json:"business_type"`
	BusinessID   string        `json:"business_id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type CollectionTask struct {
	ID          string           `json:"id"`
	Chain       Chain            `json:"chain"`
	Asset       string           `json:"asset"`
	FromAddress string           `json:"from_address"`
	ToAddress   string           `json:"to_address"`
	Amount      int64            `json:"amount"`
	Status      CollectionStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}
