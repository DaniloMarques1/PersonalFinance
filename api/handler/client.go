package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type ClientHandler struct {
	clientRepo model.IClient
	validate   *validator.Validate
}

func NewClientHandler(clientRepo model.IClient, validate *validator.Validate) *ClientHandler {
	return &ClientHandler{
		clientRepo: clientRepo,
		validate:   validate,
	}
}

func (ch *ClientHandler) SaveClient(w http.ResponseWriter, r *http.Request) {
	var clientDto dto.SaveClientRequestDto
	if err := json.NewDecoder(r.Body).Decode(&clientDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
	defer r.Body.Close()

	if err := ch.validate.Struct(clientDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}

	_, err := ch.clientRepo.FindByEmail(clientDto.Email)
	if err == nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Email already taken"})
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(clientDto.Password), bcrypt.DefaultCost)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Unnexpected error"})
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
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Unnexpected error while adding client"})
		return
	}

	util.RespondJson(w, http.StatusCreated, &dto.SaveClientResponseDto{Client: client})
}

func (ch *ClientHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var sessionDto dto.SessionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&sessionDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
        defer r.Body.Close()

	if err := ch.validate.Struct(sessionDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}

	client, err := ch.clientRepo.FindByEmail(sessionDto.Email)
	if err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Invalid email"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(client.PasswordHash, []byte(sessionDto.Password)); err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Wrong password"})
		return
	}

	token, err := util.NewToken(client.Id)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Error generating token"})
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.SessionResponseDto{Client: client, Token: token})
}
