package dto

import "github.com/danilomarques1/personalfinance/api/model"

type CreateWalletDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Total       float64 `json:"total"`
}

type GetWallets struct {
	Wallets []model.Wallet `json:"wallets"`
}

type GetWallet struct {
	Client *model.Client `json:"client"`
	Wallet *model.Wallet `json:"wallet"`
}
