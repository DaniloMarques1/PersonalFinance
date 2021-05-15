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
)

type WalletHandler struct {
	walletRepo model.IWallet
}

func NewWalletHandler(walletRepo model.IWallet) *WalletHandler {
	return &WalletHandler{
		walletRepo: walletRepo,
	}
}

func (wh *WalletHandler) SaveWallet(w http.ResponseWriter, r *http.Request) {
	var walletDto dto.CreateWalletDto
	err := json.NewDecoder(r.Body).Decode(&walletDto)
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Missing token"})
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
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Unnexpected error"})
		return
	}

	util.RespondJson(w, http.StatusCreated, &dto.CreateWalletResponse{Wallet: wallet})
}

func (wh *WalletHandler) RemoveWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, _ := strconv.Atoi(vars["wallet_id"])
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Missing token"})
		return
	}

	err = wh.walletRepo.RemoveWallet(int64(client_id), int64(wallet_id))
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: err.Error()})
		return
	}

	util.RespondJson(w, http.StatusNoContent, nil)
}

func (wh *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	client_id, err := strconv.Atoi(r.Header.Get("userId"))
	if err != nil {
		util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorDto{Message: "Missing token"})
		return
	}

	wallets, total, err := wh.walletRepo.GetWallets(int64(client_id))
	if err != nil {
		fmt.Printf("%v", err)
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Unnexpected error"})
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.GetWallets{Wallets: wallets, Total: total})
}

/*
func (wh *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    wallet_id, _ := strconv.Atoi(vars["id"])
    client_id, err := strconv.Atoi(r.Header.Get("token"))
    if err != nil {
        util.RespondJson(w, http.StatusUnauthorized, &dto.ErrorDto{Message: "Missing token"})
        return
    }

    client, wallet, err := wh.walletRepo.GetWallet(int64(wallet_id), int64(client_id))
    if err != nil {
        util.RespondJson(w, http.StatusNotFound, &dto.ErrorDto{Message: "Wallet not found"})
        return
    }

    util.RespondJson(w, http.StatusOK, &dto.GetWallet{Client: client, Wallet: wallet})
}

*/
