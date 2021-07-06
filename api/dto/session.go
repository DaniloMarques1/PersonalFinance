package dto

import "github.com/danilomarques1/personalfinance/api/model"

type SessionRequestDto struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type SessionResponseDto struct {
	Client *model.Client `json:"client"`
	Token  string        `json:"token"`
}
