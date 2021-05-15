package handler

import (
	"encoding/json"
	"fmt"
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
	// TODO
	vars := mux.Vars(r)
        wallet_id, err := strconv.Atoi(vars["wallet_id"])
        if err != nil {
                util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid wallet"})
                return
        }

        fmt.Println(wallet_id)
	var movementDto dto.AddMovementDto
        err = json.NewDecoder(r.Body).Decode(&movementDto)
	if err != nil {
		log.Fatalf("Error parsing json %v", err)
	}
        movement := model.Movement{
                Description: movementDto.Description,
                Value: movementDto.Value,
                Wallet_id: int64(wallet_id),
                Deposit: movementDto.Deposit,
        }
        err = mh.movementRepo.SaveMovement(&movement)
        if err != nil {
                util.RespondJson(w, http.StatusInternalServerError, &dto.ErrorDto{Message: err.Error()})
                return
        }

	w.WriteHeader(http.StatusOK)
}
