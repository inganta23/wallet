package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/inganta23/wallet/internal/domain"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.WalletRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetBalance(ctx context.Context, userID int64) (float64, error) {
	var balance float64
	query := "SELECT balance FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, errors.New(domain.ErrUserNotFound)
	}
	return balance, err
}

func (r *postgresRepo) WithdrawTx(ctx context.Context, userID int64, amount float64) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var currentBalance float64
	queryLock := "SELECT balance FROM users WHERE id = $1 FOR UPDATE"
	if err := tx.QueryRowContext(ctx, queryLock, userID).Scan(&currentBalance); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(domain.ErrUserNotFound)
		}
		return 0, err
	}

	if currentBalance < amount {
		return 0, errors.New(domain.ErrInsufficient)
	}

	newBalance := currentBalance - amount
	queryUpdate := "UPDATE users SET balance = $1 WHERE id = $2"
	if _, err := tx.ExecContext(ctx, queryUpdate, newBalance, userID); err != nil {
		return 0, fmt.Errorf("failed to update: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newBalance, nil
}