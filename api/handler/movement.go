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
	"github.com/gorilla/mux"
)

type MovementHandler struct {
	movementService *service.MovementService
	validate        *validator.Validate
}

func NewMovementHandler(movementService *service.MovementService, validate *validator.Validate) *MovementHandler {
	return &MovementHandler{
		movementService: movementService,
		validate:        validate,
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
	defer r.Body.Close()

	if err := mh.validate.Struct(movementDto); err != nil {
                log.Printf("Error validating struct %v\n", err)
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid Body"})
		return
	}

	movementResponse, err := mh.movementService.SaveMovement(movementDto, int64(wallet_id))
	if err != nil {
		util.HandleError(w, err)
	}

	util.RespondJson(w, http.StatusCreated, movementResponse)
}

func (mh *MovementHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wallet_id, err := strconv.Atoi(vars["wallet_id"])
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorResponseDto{Message: "Invalid wallet"})
		return
	}
	movementsResponse, err := mh.movementService.FindAll(int64(wallet_id))
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, movementsResponse)
}
