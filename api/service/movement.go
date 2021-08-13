package service

import (
	"log"
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

	// TODO check if there is a wallet for this client based on wallet_id and client_id
	/*
		TODO
		adicionar método no repositorio do movement que vai buscar uma wallet que,
		possua um wallet_id e um client_id, iguais aos recebidos via token e url param
	*/

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
		log.Printf("Error saving movement %v\n", err)
		return nil, err
	}

	movementResponse := dto.AddMovementResponseDto{Movement: movement}

	return &movementResponse, nil
}

func (ms *MovementService) FindAll(wallet_id int64) (*dto.MovementsResponseDto, error) {
	movements, err := ms.movementRepo.FindAll(wallet_id)
	if err != nil {
		return nil, err
	}

	return &dto.MovementsResponseDto{Movements: movements}, nil
}
