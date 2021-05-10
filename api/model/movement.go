package model

import "time"

type Movement struct {
	Id           int64     `json:"id"`
	Description  string    `json:"description"`
	Deposit      bool      `json:"deposit"`
	Value        float64   `json:"value"`
	MovementDate time.Time `json:"movement_date"`
	Wallet_id    int64     `json:"wallet_id,omitempty"`
}

type IMovemnt interface {
	SaveMovement(movement *Movement) error
	//GetMovements(wallet_id int64) (Wallet, []Movement, error)
}
