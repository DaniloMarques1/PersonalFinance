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
	SaveWallet(*Wallet) error
	RemoveWallet(int64, int64) error
	FindAll(int64) ([]Wallet, float64, error)
	UpdateWallet(*Wallet) error
	FindById(int64, int64) (*Wallet, error)
}
