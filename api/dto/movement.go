package dto

import "github.com/danilomarques1/personalfinance/api/model"

type AddMovementDto struct {
	Description string  `json:"description" validate:"max=100""`
	Deposit     bool    `json:"deposit" validate:"required"`
	Value       float64 `json:"value" validate:"required,gt=0"`
}

type AddMovementResponseDto struct {
	Movement model.Movement `json:"movement"`
}

type MovementsResponseDto struct {
	Movements []model.Movement `json:"movements"`
}
