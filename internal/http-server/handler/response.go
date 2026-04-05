package handler

import "github.com/google/uuid"

type ErrorResponse struct {
	Error string `json:"error"`
}

type AmountResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Amount   int       `json:"amount"`
}
