package service

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
)

type ClientService struct {
	clientRepo model.IClient
}

func NewClientService(clientRepo model.IClient) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

func (cs *ClientService) SaveClient(clientDto dto.ClientDto) (*dto.ClientDtoResponse, error) {
	_, err := cs.clientRepo.FindByEmail(clientDto.Email)
	if err == nil {
		return nil, util.NewApiError("Email already used", http.StatusBadRequest)
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(clientDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	client := model.Client{
		Id:           -1,
		Name:         clientDto.Name,
		Email:        clientDto.Email,
		PasswordHash: password_hash,
	}

	err = ch.clientRepo.SaveClient(&client)
	if err != nil {
		return nil, err
	}

	return &dto.ClientResponse{Client: client}, nil
}
