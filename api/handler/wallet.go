package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/service"
	"github.com/danilomarques1/personalfinance/api/util"
	"github.com/gorilla/mux"

	"github.com/go-playground/validator"
)

type WalletHandler struct {
	walletService *service.WalletService
	validate      *validator.Validate
}

func NewWalletHandler(walletService *service.WalletService, validate *validator.Validate) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		validate:      validate,
	}
}

func (wh *WalletHandler) SaveWallet(w http.ResponseWriter, r *http.Request) {
	var walletDto dto.SaveWalletRequestDto
	if err := json.NewDecoder(r.Body).Decode(&walletDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
	defer r.Body.Close()

	if err := wh.validate.Struct(walletDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}

	clientId, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}

	walletResponse, err := wh.walletService.SaveWallet(walletDto, int64(clientId))
	if err != nil {
		util.HandleError(w, err)
	}

	util.RespondJson(w, http.StatusCreated, walletResponse)
}

func (wh *WalletHandler) RemoveWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, "Invalid params")
		return
	}
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}

	err = wh.walletService.RemoveWallet(int64(wallet_id), int64(client_id))
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusNoContent, nil)
}

func (wh *WalletHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}
	defer r.Body.Close()

	walletsResponse, err := wh.walletService.FindAll(int64(client_id))
	if err != nil {
		log.Printf("Error searching wallets %v", err)
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, walletsResponse)
}

func (wh *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}

	vars := mux.Vars(r)
	walletId, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid params"})
		return
	}

	var walletUpdate dto.WalletUpdateRequestDto
	if err := json.NewDecoder(r.Body).Decode(&walletUpdate); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}
	defer r.Body.Close()

	if err := wh.validate.Struct(walletUpdate); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid body"})
		return
	}

	err = wh.walletService.UpdateWallet(walletUpdate, int64(walletId), int64(clientId))
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusNoContent, nil)
}
