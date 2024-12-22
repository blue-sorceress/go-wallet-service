package service

import (
	"go-wallet-service/internal/model"
	"go-wallet-service/internal/repository"
)

type WalletService struct {
	walletRepository repository.WalletRepository
}

func NewWalletService(walletRepository repository.WalletRepository) *WalletService {
	return &WalletService{walletRepository: walletRepository}
}

func (ws *WalletService) GetWalletByUserId(userId int) ([]model.Wallet, error) {
	return ws.walletRepository.GetWalletByUserId(userId)
}

func (ws *WalletService) GetById(walletId int) (model.Wallet, error) {
	return ws.walletRepository.GetById(walletId)
}

func (ws *WalletService) Update(wallet model.Wallet, amount float64, transactionType string) error {
	return ws.walletRepository.Update(wallet, amount, transactionType)
}
