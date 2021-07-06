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

type IMovement interface {
	SaveMovement(movement *Movement) error
	FindAll(wallet_id int64) ([]Movement, error)
	CanWithDraw(wallet_id int64, value float64) (bool, error)
}
