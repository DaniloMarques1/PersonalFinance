package dto

import (
	"github.com/danilomarques1/personalfinance/api/model"
)

type SaveClientRequestDto struct {
	Name     string `json:"name" validate:"required,max=60"`
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=6,max=20"` // TODO verificar tamanho (len) de um hash com password de 20
}

type SaveClientResponseDto struct {
	Client model.Client `json:"client"`
}
