package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"

	"golang.org/x/crypto/bcrypt"
)

type ClientHandler struct {
	clientRepo model.IClient
}

func NewClientHandler(clientRepo model.IClient) *ClientHandler {
	return &ClientHandler{
		clientRepo: clientRepo,
	}
}

func (ch *ClientHandler) SaveClient(w http.ResponseWriter, r *http.Request) {
	var clientDto dto.ClientDto
	if err := json.NewDecoder(r.Body).Decode(&clientDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}

	if clientDto.Name == "" || clientDto.Email == "" || clientDto.Password == "" {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}

	_, err := ch.clientRepo.FindByEmail(clientDto.Email)
	if err == nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Email already taken"})
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(clientDto.Password), bcrypt.DefaultCost)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Unnexpected error"})
		return
	}

	client := model.Client{
		Id:           -1,
		Name:         clientDto.Name,
		Email:        clientDto.Email,
		PasswordHash: password_hash,
	}

	err = ch.clientRepo.SaveClient(&client)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Unnexpected error while adding client"})
		return
	}

	util.RespondJson(w, http.StatusCreated, &client)
}

func (ch *ClientHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var sessionDto dto.SessionDto
	err := json.NewDecoder(r.Body).Decode(&sessionDto)
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}

	client, err := ch.clientRepo.FindByEmail(sessionDto.Email)
	if err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorDto{Message: "Invalid email"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(client.PasswordHash, []byte(sessionDto.Password)); err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorDto{Message: "Wrong password"})
		return
	}

	token, err := util.NewToken(client.Id)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Error generating token"})
		fmt.Printf("%v\n", err)
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.SessionResponseDto{Client: client, Token: token})
}
