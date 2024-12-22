package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"go-wallet-service/internal/model"
)

type UserRepository interface {
	GetById(userId int) (model.User, error)
	GetByIds(userId []int) ([]model.User, error)
	GetUserByToken(token string) (int, error)
	GetUserByUserId(userId int) ([]model.User, error)
	GetUserTransactionsByUserId(userId int) ([]model.Transaction, error)
	Transfer(senderUserId int, senderWalletId int, senderWalletBalance float64, receiverUserId int, receiverWalletId int, receiverWalletBalance float64, amount float64) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetById(userId int) (model.User, error) {
	var user model.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE id = $1", userId)

	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) GetByIds(userIds []int) ([]model.User, error) {
	var users []model.User

	strIds := make([]string, len(userIds))
	for i, id := range userIds {
		strIds[i] = fmt.Sprintf("%d", id)
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id IN (%s)", strings.Join(strIds, ","))
	err := ur.db.Select(&users, query)

	if err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (ur *userRepository) GetUserByUserId(userId int) ([]model.User, error) {
	var user []model.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetUserByToken(token string) (int, error) {
	var userId int
	err := ur.db.Get(&userId, "SELECT user_id FROM oauth WHERE token = $1", token)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (ur *userRepository) GetUserTransactionsByUserId(userId int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := ur.db.Select(&transactions, "SELECT * FROM transactions WHERE sender_user_id = $1 ORDER BY creation_date DESC", userId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *userRepository) Transfer(senderUserId int, senderWalletId int, senderWalletBalance float64, receiverUserId int, receiverWalletId int, receiverWalletBalance float64, amount float64) error {
	_, err := r.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", senderWalletBalance, senderWalletId)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", receiverWalletBalance, receiverWalletId)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO transactions (sender_user_id, sender_wallet_id, receiver_user_id, receiver_wallet_id, amount, type) VALUES ($1, $2, $3, $4, $5, $6)", senderUserId, senderWalletId, receiverUserId, receiverWalletId, amount, "transfer")
	if err != nil {
		return err
	}

	return nil
}
