package service

import (
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
)

type MovementService struct {
	movementRepo model.IMovement
}

func NewMovementService(movementRepo model.IMovement) *MovementService {
	return &MovementService{
		movementRepo: movementRepo,
	}
}

func (ms *MovementService) SaveMovement(movementDto dto.AddMovementDto, wallet_id int64) (*dto.AddMovementResponseDto, error) {
	if !movementDto.Deposit {
		canWithDraw, err := ms.movementRepo.CanWithDraw(wallet_id, movementDto.Value)
		if err != nil {
			return nil, err
		}
		if !canWithDraw {
			return nil, util.NewApiError("You do not have enough to withdraw", http.StatusUnauthorized)
		}
	}

	movement := model.Movement{
		Description: movementDto.Description,
		Value:       movementDto.Value,
		Deposit:     movementDto.Deposit,
		Wallet_id:   wallet_id,
	}
	err := ms.movementRepo.SaveMovement(&movement)
	if err != nil {
		return nil, err
	}

	movementResponse := dto.AddMovementResponseDto{Movement: movement}

	return &movementResponse, nil
}
