package repository

import (
	"github.com/jmoiron/sqlx"

	"go-wallet-service/internal/model"
)

type WalletRepository interface {
	GetWalletByUserId(userId int) ([]model.Wallet, error)
	GetById(walletId int) (model.Wallet, error)
	Update(wallet model.Wallet, amount float64, transactionType string) error
}

type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *walletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetWalletByUserId(userId int) ([]model.Wallet, error) {
	var wallet []model.Wallet
	err := r.db.Select(&wallet, "SELECT * FROM wallets WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) GetById(walletId int) (model.Wallet, error) {
	var wallet model.Wallet

	err := r.db.Get(&wallet, "SELECT * FROM wallets WHERE id = $1", walletId)
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}

func (r *walletRepository) Update(wallet model.Wallet, amount float64, transactionType string) error {
	_, err := r.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", wallet.Balance, wallet.ID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO transactions (sender_user_id, sender_wallet_id, receiver_user_id, receiver_wallet_id, amount, type) VALUES ($1, $2, $3, $4, $5, $6)", wallet.UserId, wallet.ID, wallet.UserId, wallet.ID, amount, transactionType)
	return err
}
