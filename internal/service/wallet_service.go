package service

import (
	"context"
	"errors"

	"github.com/inganta23/wallet/internal/domain"
)

type walletService struct {
	repo domain.WalletRepository
}

func NewWalletService(repo domain.WalletRepository) domain.WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetBalance(ctx context.Context, userID int64) (float64, error) {
	return s.repo.GetBalance(ctx, userID)
}

func (s *walletService) Withdraw(ctx context.Context, userID int64, amount float64) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}
	return s.repo.WithdrawTx(ctx, userID, amount)
}