package model

import "time"

type Wallet struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Total        float64   `json:"total"`
	Created_date time.Time `json:"created_date"`
	Client_id    int64     `json:"client_id"`
}

type IWallet interface {
	SaveWallet(wallet *Wallet) error
	RemoveWallet(client_id, wallet_id int64) error
	//GetWallets(client_id int64) ([]Wallet, float64, error)
	FindAll(client_id int64) ([]Wallet, float64, error)
}
