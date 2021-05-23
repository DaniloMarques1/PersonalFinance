package dto

import (
	"github.com/danilomarques1/personalfinance/api/model"
)

type ClientDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type ClientResponseDto struct {
	Client model.Client `json:"client"`
}
