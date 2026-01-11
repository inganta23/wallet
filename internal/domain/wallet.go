package domain

import "context"

var (
	ErrUserNotFound = "user not found"
	ErrInsufficient = "insufficient funds"
)

type User struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

type WalletRepository interface {
	GetBalance(ctx context.Context, userID int64) (float64, error)
	WithdrawTx(ctx context.Context, userID int64, amount float64) (float64, error)
}

type WalletService interface {
	GetBalance(ctx context.Context, userID int64) (float64, error)
	Withdraw(ctx context.Context, userID int64, amount float64) (float64, error)
}