package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/inganta23/wallet/internal/domain"
	"github.com/inganta23/wallet/internal/service"
)

type MockRepository struct {
	balanceToReturn float64
	withdrawError   error
	getBalanceError error
}

func (m *MockRepository) GetBalance(ctx context.Context, userID int64) (float64, error) {
	if m.getBalanceError != nil {
		return 0, m.getBalanceError
	}
	return m.balanceToReturn, nil
}

func (m *MockRepository) WithdrawTx(ctx context.Context, userID int64, amount float64) (float64, error) {
	if m.withdrawError != nil {
		return 0, m.withdrawError
	}
	return m.balanceToReturn - amount, nil
}

func TestWithdraw(t *testing.T) {
	// Define Test Cases
	tests := []struct {
		name          string
		inputAmount   float64
		inputUserID   int64
		mockBalance   float64
		mockError     error
		expectedError string
		expectSuccess bool
	}{
		{
			name:          "Success: Valid Withdrawal",
			inputAmount:   50.00,
			inputUserID:   1,
			mockBalance:   100.00,
			mockError:     nil,
			expectSuccess: true,
		},
		{
			name:          "Fail: Negative Amount (Logic Check)",
			inputAmount:   -10.00, // Invalid input
			inputUserID:   1,
			mockBalance:   100.00,
			mockError:     nil,
			expectedError: "amount must be positive",
			expectSuccess: false,
		},
		{
			name:          "Fail: Insufficient Funds (DB Logic)",
			inputAmount:   200.00,
			inputUserID:   1,
			mockBalance:   100.00,
			mockError:     errors.New(domain.ErrInsufficient),
			expectedError: domain.ErrInsufficient,
			expectSuccess: false,
		},
		{
			name:          "Fail: User Not Found",
			inputAmount:   50.00,
			inputUserID:   999,
			mockBalance:   0,
			mockError:     errors.New(domain.ErrUserNotFound),
			expectedError: domain.ErrUserNotFound,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				balanceToReturn: tt.mockBalance,
				withdrawError:   tt.mockError,
			}

			svc := service.NewWalletService(mockRepo)

			newBalance, err := svc.Withdraw(context.Background(), tt.inputUserID, tt.inputAmount)

			if tt.expectSuccess {
				if err != nil {
					t.Fatalf("expected success, got error: %v", err)
				}
				expectedBalance := tt.mockBalance - tt.inputAmount
				if newBalance != expectedBalance {
					t.Errorf("expected balance %.2f, got %.2f", expectedBalance, newBalance)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error msg '%s', got '%s'", tt.expectedError, err.Error())
				}
			}
		})
	}
}