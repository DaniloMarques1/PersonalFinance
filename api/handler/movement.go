package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
	"github.com/gorilla/mux"
)

type MovementHandler struct {
	movementRepo model.IMovement
}

func NewMovementHandler(movementRepo model.IMovement) *MovementHandler {
	return &MovementHandler{
		movementRepo: movementRepo,
	}
}

func (mh *MovementHandler) SaveMovement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid wallet"})
		return
	}

	var movementDto dto.AddMovementDto
	err = json.NewDecoder(r.Body).Decode(&movementDto)
	if err != nil {
		log.Fatalf("Error parsing json %v", err)
                util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid Body"})
                return
	}
	movement := model.Movement{
		Description: movementDto.Description,
		Value:       movementDto.Value,
		Wallet_id:   int64(wallet_id),
		Deposit:     movementDto.Deposit,
	}
	err = mh.movementRepo.SaveMovement(&movement)
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: err.Error()})
		return
	}

	util.RespondJson(w, http.StatusCreated, &dto.AddMovementResponseDto{Movement: movement})
}

func (mh *MovementHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid wallet"})
		return
	}
	movements, err := mh.movementRepo.GetMovements(int64(wallet_id))
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: "Unexpected error"})
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.GetMovements{Movements: movements})
}
