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

func (ms *MovementService) SaveMovement(movementDto dto.AddMovementDto, walletId, clientId int64) (*dto.AddMovementResponseDto, error) {
	_, err := ms.movementRepo.FindMovementWallet(walletId, clientId) // if returns error, means it did not find any wallet
	if err != nil {
		return nil, util.NewApiError("You do not have a wallet with that id", http.StatusNotFound)
	}

	if !movementDto.Deposit {
		canWithDraw, err := ms.movementRepo.CanWithDraw(walletId, movementDto.Value)
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
		Wallet_id:   walletId,
	}
	err = ms.movementRepo.SaveMovement(&movement)
	if err != nil {
		log.Printf("Error saving movement %v\n", err)
		return nil, err
	}

	movementResponse := dto.AddMovementResponseDto{Movement: movement}

	return &movementResponse, nil
}

func (ms *MovementService) FindAll(walletId int64) (*dto.MovementsResponseDto, error) {
	movements, err := ms.movementRepo.FindAll(walletId)
	if err != nil {
		return nil, err
	}

	return &dto.MovementsResponseDto{Movements: movements}, nil
}
