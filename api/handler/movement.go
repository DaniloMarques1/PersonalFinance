package handler

import (
        "net/http"

        "github.com/danilomarques1/personalfinance/api/model"
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
}
