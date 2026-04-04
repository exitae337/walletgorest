package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletRepository struct {
	storage *Storage
}

func NewWalletRepo(storage *Storage) *WalletRepository {
	return &WalletRepository{storage: storage}
}

func (w *WalletRepository) GetAmountOfMoney(ctx context.Context, id uuid.UUID) (int, error) {
	const op = "postgres.repo.GetAmountOfMoney"

	query := "SELECT amount FROM walletdb WHERE walletid = $1"
	var amount int

	err := w.storage.Pool().QueryRow(ctx, query, id).Scan(&amount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("%s: wallet not found: %s", op, id)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if amount < 0 {
		return 0, fmt.Errorf("%s: negative balance detected: %d", op, amount)
	}

	return amount, nil
}

func (w *WalletRepository) WithdrawMoney(ctx context.Context, amount int, id uuid.UUID) error {
	const op = "postgres.repo.WithdrawMoney"

	if amount <= 0 {
		return fmt.Errorf("%s: amount must be positive, got %d", op, amount)
	}

	tx, err := w.storage.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx, `
        UPDATE walletdb 
        SET amount = amount - $1 
        WHERE walletid = $2 AND amount >= $1
    `, amount, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		var exists bool
		err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM walletdb WHERE walletid = $1)", id).Scan(&exists)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: wallet not found: %s", op, id)
		}
		return fmt.Errorf("%s: insufficient funds", op)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (w *WalletRepository) DepositMoney(ctx context.Context, amount int, id uuid.UUID) error {
	const op = "postgres.repo.DepositMoney"

	if amount <= 0 {
		return fmt.Errorf("%s: amount must be positive, got %d", op, amount)
	}

	tx, err := w.storage.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx, `
        UPDATE walletdb 
        SET amount = amount + $1 
        WHERE walletid = $2
    `, amount, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: wallet not found: %s", op, id)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (w *WalletRepository) CreateWallet(ctx context.Context, id uuid.UUID, initialAmount int) error {
	const op = "postgres.repo.CreateWallet"

	if initialAmount < 0 {
		return fmt.Errorf("%s: initial amount cannot be negative", op)
	}

	_, err := w.storage.Pool().Exec(ctx,
		"INSERT INTO walletdb (walletid, amount) VALUES ($1, $2)",
		id, initialAmount)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
