package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/service"
	"github.com/danilomarques1/personalfinance/api/util"

	"github.com/go-playground/validator"
)

type ClientHandler struct {
	clientService *service.ClientService
	validate      *validator.Validate
}

func NewClientHandler(clientService *service.ClientService, validate *validator.Validate) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		validate:      validate,
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

	clientResponse, err := ch.clientService.SaveClient(clientDto)
	if err != nil {
		log.Printf("Error saving a client %v", err)
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusCreated, clientResponse)
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

	sessionResponse, err := ch.clientService.CreateSession(sessionDto)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, sessionResponse)
}

func (ch *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var updateClientDto dto.UpdateClientRequestDto
	if err := json.NewDecoder(r.Body).Decode(&updateClientDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
	defer r.Body.Close()

	if err := ch.validate.Struct(updateClientDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
	clientId, err := strconv.Atoi(r.Header.Get("client_id"))
	if err != nil {
		//
	}
	err = ch.clientService.UpdateClient(int64(clientId), updateClientDto)
	if err != nil {
		log.Printf("error %v\n", err)
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
