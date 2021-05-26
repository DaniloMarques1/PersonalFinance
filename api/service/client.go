package service

import (
        "database/sql"
        "fmt"
        "net/http"

        "github.com/danilomarques1/personalfinance/api/model"
        "github.com/danilomarques1/personalfinance/api/dto"
)

type ClientService struct {
        clientRepo model.IClient
}

type UnauthorizedError struct {
        error
}

func NewClientService(clientRepo model.IClient) *ClientService {
        return &ClientService{
                clientRepo: clientRepo,
        }
}

func (cs *ClientService) SaveClient(clientDto dto.ClientDto) (*dto.ClientDtoResponse, error) {
	_, err := ch.clientRepo.FindByEmail(clientDto.Email)
	if err == nil {
                // TODO return a different error becasue this one is nil
		return nil, err
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(clientDto.Password), bcrypt.DefaultCost)
	if err != nil {
                // TODO format error message
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
                // TODO format error message
		return nil, err
	}

        return &dto.ClientResponse{Client: client}, nil
}
