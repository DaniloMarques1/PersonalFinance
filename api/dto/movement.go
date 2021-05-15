package dto

import "github.com/danilomarques1/personalfinance/api/model"

type AddMovementDto struct {
	Description string  `json:"description"`
	Deposit     bool    `json:"deposit"`
	Value       float64 `json:"value"`
}

type AddMovementResponseDto struct {
	Movement model.Movement `json:"movement"`
}

type GetMovements struct {
	Movements []model.Movement `json:"movements"`
}
