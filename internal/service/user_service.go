package service

import (
	"go-wallet-service/internal/model"
	"go-wallet-service/internal/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetById(userId int) (model.User, error) {
	return us.userRepository.GetById(userId)
}

func (us *UserService) GetByIds(userIds []int) ([]model.User, error) {
	return us.userRepository.GetByIds(userIds)
}

func (us *UserService) GetUserTransactionsByUserId(userId int) ([]model.Transaction, error) {
	return us.userRepository.GetUserTransactionsByUserId(userId)
}

func (us *UserService) GetUserByToken(token string) (int, error) {
	return us.userRepository.GetUserByToken(token)
}

func (us *UserService) Transfer(senderUserId int, senderWalletId int, senderWalletBalance float64, receiverUserId int, receiverWalletId int, receiverWalletBalance float64, amount float64) error {
	return us.userRepository.Transfer(senderUserId, senderWalletId, senderWalletBalance, receiverUserId, receiverWalletId, receiverWalletBalance, amount)
}
