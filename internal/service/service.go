package service

import (
	"context"

	"github.com/exitae337/walletgorest/internal/storage/postgres"
	"github.com/google/uuid"
)

type WalletService struct {
	repo postgres.WalletRepo
}

func NewWalletService(repo postgres.WalletRepo) *WalletService {
	return &WalletService{repo: repo}
}

func (w *WalletService) CreateNewWallet(ctx context.Context, id uuid.UUID, amount int) error {
	return w.repo.CreateWallet(ctx, id, amount)
}

func (w *WalletService) Deposit(ctx context.Context, amount int, id uuid.UUID) error {
	return w.repo.DepositMoney(ctx, amount, id)
}

func (w *WalletService) Withdraw(ctx context.Context, amount int, id uuid.UUID) error {
	return w.repo.WithdrawMoney(ctx, amount, id)
}

func (w *WalletService) GetAmount(ctx context.Context, id uuid.UUID) (int, error) {
	return w.repo.GetAmountOfMoney(ctx, id)
}
