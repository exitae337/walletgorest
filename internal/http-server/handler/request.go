package handler

import "github.com/google/uuid"

type OperationRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"` // DEPOSIT or WITHDRAW
	Amount        int       `json:"amount"`
}

type CreateWalletRequest struct {
	InitialAmount int `json:"initialAmount,omitempty"`
}
