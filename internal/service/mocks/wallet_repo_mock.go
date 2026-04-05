package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) GetAmountOfMoney(ctx context.Context, id uuid.UUID) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockWalletRepo) WithdrawMoney(ctx context.Context, amount int, id uuid.UUID) error {
	args := m.Called(ctx, amount, id)
	return args.Error(0)
}

func (m *MockWalletRepo) DepositMoney(ctx context.Context, amount int, id uuid.UUID) error {
	args := m.Called(ctx, amount, id)
	return args.Error(0)
}

func (m *MockWalletRepo) CreateWallet(ctx context.Context, id uuid.UUID, initialAmount int) error {
	args := m.Called(ctx, id, initialAmount)
	return args.Error(0)
}
