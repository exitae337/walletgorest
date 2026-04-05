package handler

import (
	"encoding/json"
	"net/http"

	"github.com/exitae337/walletgorest/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WalletHandler struct {
	serv *service.WalletService
}

func NewWalletHandler(serv *service.WalletService) *WalletHandler {
	return &WalletHandler{
		serv: serv,
	}
}

func (h *WalletHandler) RegisterRoutes(router chi.Router) {
	router.Route("/", func(r chi.Router) {
		r.Post("/wallets", h.CreateWallet)
		r.Post("/wallet", h.OpearationWallet)
		r.Get("/wallets/{wallet_uid}", h.GetAmount)
	})
}

// Create wallet func
func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	walletID := uuid.New()

	err := h.serv.CreateNewWallet(r.Context(), walletID, req.InitialAmount)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"walletId": walletID,
		"message":  "wallet created without errors",
	})
}

// Operation with wallet func
func (h *WalletHandler) OpearationWallet(w http.ResponseWriter, r *http.Request) {
	var req OperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletID == uuid.Nil {
		respondWithError(w, "walletId cannot be nil", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		respondWithError(w, "walletId cannot be nil", http.StatusBadRequest)
		return
	}

	var err error

	switch req.OperationType {
	case "DEPOSIT":
		err = h.serv.Deposit(r.Context(), req.Amount, req.WalletID)
	case "WITHDRAW":
		err = h.serv.Withdraw(r.Context(), req.Amount, req.WalletID)
	default:
		respondWithError(w, "invalid operation type: must be DEPOSIT or WITHDRAW", http.StatusBadRequest)
		return
	}

	if err != nil {
		switch err.Error() {
		case "insufficient funds":
			respondWithError(w, err.Error(), http.StatusBadRequest)
		case "wallet not found":
			respondWithError(w, err.Error(), http.StatusNotFound)
		default:
			respondWithError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "operation completed successfully",
	})
}

// Get amount func
func (h *WalletHandler) GetAmount(w http.ResponseWriter, r *http.Request) {
	walletUUIDStr := chi.URLParam(r, "wallet_uid")

	if walletUUIDStr == "" {
		respondWithError(w, "wallet_uid param can't be empty", http.StatusBadRequest)
		return
	}

	walletID, err := uuid.Parse(walletUUIDStr)
	if err != nil {
		respondWithError(w, "failed to parse UUID", http.StatusBadRequest)
		return
	}

	amount, err := h.serv.GetAmount(r.Context(), walletID)
	if err != nil {
		if err.Error() == "wallet not found" {
			respondWithError(w, err.Error(), http.StatusNotFound)
			return
		}
		respondWithError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, AmountResponse{
		WalletID: walletID,
		Amount:   amount,
	})
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, message string, status int) {
	respondWithJSON(w, status, ErrorResponse{Error: message})
}
