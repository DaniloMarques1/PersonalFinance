package dto

import "github.com/danilomarques1/personalfinance/api/model"

type SessionDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SessionResponseDto struct {
	Client *model.Client `json:"client"`
	Token  string        `json:"token"`
}
