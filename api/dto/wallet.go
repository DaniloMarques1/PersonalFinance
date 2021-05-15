package dto

import "github.com/danilomarques1/personalfinance/api/model"

type CreateWalletDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Total       float64 `json:"total"`
}

type CreateWalletResponse struct {
	Wallet model.Wallet `json:"wallet"`
}

type GetWallets struct {
	Wallets []model.Wallet `json:"wallets"`
	Total   float64        `json:"total"`
}

type GetWallet struct {
	//Wallet *model.Wallet `json:"wallet"`
	Movements []model.Movement `json:"movements"`
}
