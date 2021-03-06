package service

import (
	"log"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"

	"golang.org/x/crypto/bcrypt"
)

type ClientService struct {
	clientRepo model.IClient
}

func NewClientService(clientRepo model.IClient) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

func (cs *ClientService) SaveClient(clientDto dto.SaveClientRequestDto) (*dto.SaveClientResponseDto, error) {
	_, err := cs.clientRepo.FindByEmail(clientDto.Email)
	if err == nil {
		log.Printf("Same email was already found\n")
		return nil, util.NewApiError("Email already used", http.StatusBadRequest)
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(clientDto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating a password hash %v", err)
		return nil, err
	}

	client := model.Client{
		Id:           -1,
		Name:         clientDto.Name,
		Email:        clientDto.Email,
		PasswordHash: password_hash,
	}

	err = cs.clientRepo.SaveClient(&client)
	if err != nil {
		log.Printf("Error saving client in the data base %v", err)
		return nil, err
	}

	return &dto.SaveClientResponseDto{Client: client}, nil
}

func (cs *ClientService) CreateSession(sessionDto dto.SessionRequestDto) (*dto.SessionResponseDto, error) {
	client, err := cs.clientRepo.FindByEmail(sessionDto.Email)
	if err != nil {
		return nil, util.NewApiError("Invalid email", http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword(client.PasswordHash, []byte(sessionDto.Password)); err != nil {
		return nil, util.NewApiError("Wrong password", http.StatusUnauthorized)
	}

	token, refreshToken, err := util.NewToken(client.Id)
	if err != nil {
		return nil, err
	}

	return &dto.SessionResponseDto{
		Client:       client,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (cs *ClientService) RefreshSession(clientId int64) (*dto.SessionResponseDto, error) {
	client, err := cs.clientRepo.FindById(clientId)
	if err != nil {
		return nil, err
	}

	token, refreshToken, err := util.NewToken(clientId)
	if err != nil {
		return nil, err
	}
	session := dto.SessionResponseDto{
		Client:       client,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return &session, nil
}

func (cs *ClientService) UpdateClient(clientId int64, updateClientDto dto.UpdateClientRequestDto) error {
	if updateClientDto.Password != updateClientDto.ConfirmPassword {
		return util.NewApiError("Password and confirm password does not match", http.StatusBadRequest)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(updateClientDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	client := model.Client{
		Id:           clientId,
		Name:         updateClientDto.Name,
		Email:        updateClientDto.Email,
		PasswordHash: passwordHash,
	}

	err = cs.clientRepo.UpdateClient(&client)
	if err != nil {
		return err
	}

	return nil
}
