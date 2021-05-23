package dto

import (
	"github.com/danilomarques1/personalfinance/api/model"
)

type SaveClientRequestDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SaveClientResponseDto struct {
	Client model.Client `json:"client"`
}
