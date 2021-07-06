package service

import (
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

func (ws *WalletService) SaveWallet(walletDto dto.SaveWalletRequestDto, clientId int64) (*dto.SaveWalletResponseDto, error) {
	wallet := model.Wallet{
		Id:          -1,
		Name:        walletDto.Name,
		Description: walletDto.Description,
		Client_id:   clientId,
	}

        err := ws.walletRepo.SaveWallet(&wallet)
	if err != nil {
		return nil, err
	}

	return &dto.SaveWalletResponseDto{Wallet: wallet}, nil
}
