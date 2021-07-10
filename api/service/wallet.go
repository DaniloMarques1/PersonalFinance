package service

import (
        "log"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
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
