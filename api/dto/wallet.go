package dto

import "github.com/danilomarques1/personalfinance/api/model"

type SaveWalletRequestDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"max=150"`
}

type SaveWalletResponseDto struct {
	Wallet model.Wallet `json:"wallet"`
}

type WalletsResponseDto struct {
	Wallets []model.Wallet `json:"wallets"`
	Total   float64        `json:"total"`
}

type WalletResponseDto struct {
	Movements []model.Movement `json:"movements"`
}
