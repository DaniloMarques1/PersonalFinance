package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
	"github.com/gorilla/mux"

	"github.com/go-playground/validator"
)

type WalletHandler struct {
	walletRepo model.IWallet
	validate   *validator.Validate
}

func NewWalletHandler(walletRepo model.IWallet, validate *validator.Validate) *WalletHandler {
	return &WalletHandler{
		walletRepo: walletRepo,
		validate:   validate,
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

	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}

	wallet := model.Wallet{
		Id:          -1,
		Name:        walletDto.Name,
		Description: walletDto.Description,
		Client_id:   int64(client_id),
	}

	err = wh.walletRepo.SaveWallet(&wallet)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Unnexpected error"})
		return
	}

	util.RespondJson(w, http.StatusCreated, &dto.SaveWalletResponseDto{Wallet: wallet})
}

func (wh *WalletHandler) RemoveWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, _ := strconv.Atoi(vars["wallet_id"])
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Missing token"})
		return
	}

	err = wh.walletRepo.RemoveWallet(int64(client_id), int64(wallet_id))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: err.Error()})
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

	wallets, total, err := wh.walletRepo.FindAll(int64(client_id))
	if err != nil {
		fmt.Printf("%v", err)
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Unnexpected error"})
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.WalletsResponseDto{Wallets: wallets, Total: total})
}
