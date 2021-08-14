package service

import (
	"log"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
)

type WalletService struct {
	walletRepo model.IWallet
}

func NewWalletService(walletRepo model.IWallet) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

func (ws *WalletService) SaveWallet(walletDto dto.SaveWalletRequestDto, client_id int64) (*dto.SaveWalletResponseDto, error) {
	wallet := model.Wallet{
		Id:          -1,
		Name:        walletDto.Name,
		Description: walletDto.Description,
		Client_id:   client_id,
	}

	err := ws.walletRepo.SaveWallet(&wallet)
	if err != nil {
		return nil, err
	}

	return &dto.SaveWalletResponseDto{Wallet: wallet}, nil
}

func (ws *WalletService) RemoveWallet(wallet_id, client_id int64) error {
	err := ws.walletRepo.RemoveWallet(wallet_id, client_id)
	if err != nil {
		log.Printf("Error calling repository %v", err)
		return err
	}

	return nil
}

func (ws *WalletService) FindAll(client_id int64) (*dto.WalletsResponseDto, error) {
	wallets, total, err := ws.walletRepo.FindAll(client_id)
	if err != nil {
		return nil, err
	}

	return &dto.WalletsResponseDto{Wallets: wallets, Total: total}, nil
}

func (ws *WalletService) UpdateWallet(walletUpdate dto.WalletUpdateRequestDto, walletId, clientId int64) error {
	_, err := ws.walletRepo.FindById(walletId, clientId)
	if err != nil {
		return util.NewApiError("Wallet not found", http.StatusNotFound)
	}
	wallet := model.Wallet{
		Id:          walletId,
		Name:        walletUpdate.Name,
		Description: walletUpdate.Description,
		Client_id:   clientId,
	}
	err = ws.walletRepo.UpdateWallet(&wallet)
	if err != nil {
		return err
	}

	return nil
}
