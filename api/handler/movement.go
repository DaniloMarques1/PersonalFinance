package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type MovementHandler struct {
	movementRepo model.IMovement
	validate     *validator.Validate
}

func NewMovementHandler(movementRepo model.IMovement, validate *validator.Validate) *MovementHandler {
	return &MovementHandler{
		movementRepo: movementRepo,
		validate:     validate,
	}
}

func (mh *MovementHandler) SaveMovement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid wallet"})
		return
	}

	var movementDto dto.AddMovementDto
	if err = json.NewDecoder(r.Body).Decode(&movementDto); err != nil {
		log.Fatalf("Error parsing json %v", err)
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid Body"})
		return
	}

	if err := mh.validate.Struct(movementDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid Body"})
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
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: err.Error()})
		return
	}

	util.RespondJson(w, http.StatusCreated, &dto.AddMovementResponseDto{Movement: movement})
}

func (mh *MovementHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid wallet"})
		return
	}
	movements, err := mh.movementRepo.FindAll(int64(wallet_id))
	if err != nil {
		util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorResponseDto{Message: "Unexpected error"})
		return
	}

	util.RespondJson(w, http.StatusOK, &dto.MovementsResponseDto{Movements: movements})
}
