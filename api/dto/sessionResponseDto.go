package dto

import "github.com/danilomarques1/personalfinance/api/model"

type SessionResponseDto struct {
	Client *model.Client `json:"client"`
	Token  string        `json:"token"`
}
