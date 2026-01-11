package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/inganta23/wallet/internal/domain"
)

type WalletHandler struct {
	Service domain.WalletService
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid or missing user_id")
		return
	}

	balance, err := h.Service.GetBalance(r.Context(), userID)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
		"balance": balance,
	})
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newBalance, err := h.Service.Withdraw(r.Context(), req.UserID, req.Amount)
	if err != nil {
		slog.Error("Withdraw failed", "user_id", req.UserID, "error", err)
		handleDomainError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"status":      "success",
		"new_balance": newBalance,
	})
}


func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func handleDomainError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case domain.ErrInsufficient:
		respondError(w, http.StatusBadRequest, "Insufficient funds")
	case domain.ErrUserNotFound:
		respondError(w, http.StatusNotFound, "User not found")
	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}