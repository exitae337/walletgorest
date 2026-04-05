package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exitae337/walletgorest/internal/service"
	"github.com/exitae337/walletgorest/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletService_CreateNewWallet_Success(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletId := uuid.New()
	amount := 1000

	mockRepo.On("CreateWallet", ctx, walletId, amount).Return(nil)

	err := service.CreateNewWallet(ctx, walletId, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_CreateNewWallet_Error(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := 1000
	expectedErr := errors.New("database error")

	mockRepo.On("CreateWallet", ctx, walletID, amount).Return(expectedErr)

	err := service.CreateNewWallet(ctx, walletID, amount)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_Deposit_Success(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := 500

	mockRepo.On("DepositMoney", ctx, amount, walletID).Return(nil)

	err := service.Deposit(ctx, amount, walletID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_Deposit_Error(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := 500
	expectedErr := errors.New("wallet not found")

	mockRepo.On("DepositMoney", ctx, amount, walletID).Return(expectedErr)

	err := service.Deposit(ctx, amount, walletID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_Withdraw_InsufficientFunds(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := 10000
	expectedErr := errors.New("insufficient funds")

	mockRepo.On("WithdrawMoney", ctx, amount, walletID).Return(expectedErr)

	err := service.Withdraw(ctx, amount, walletID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_Withdraw_WalletNotFound(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := 100
	expectedErr := errors.New("wallet not found")

	mockRepo.On("WithdrawMoney", ctx, amount, walletID).Return(expectedErr)

	err := service.Withdraw(ctx, amount, walletID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_GetAmount_Success(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	expectedAmount := 1500

	mockRepo.On("GetAmountOfMoney", ctx, walletID).Return(expectedAmount, nil)

	amount, err := service.GetAmount(ctx, walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAmount, amount)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_GetAmount_WalletNotFound(t *testing.T) {
	mockRepo := new(mocks.MockWalletRepo)
	service := service.NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	expectedErr := errors.New("wallet not found")

	mockRepo.On("GetAmountOfMoney", ctx, walletID).Return(0, expectedErr)

	amount, err := service.GetAmount(ctx, walletID)

	assert.Error(t, err)
	assert.Equal(t, 0, amount)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
