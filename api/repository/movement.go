package repository

import (
	"database/sql"
	"log"

	"github.com/danilomarques1/personalfinance/api/model"
)

type MovementRepository struct {
	db *sql.DB
}

func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{
		db: db,
	}
}

func (mr *MovementRepository) SaveMovement(movement *model.Movement) error {
	stmt, err := mr.db.Prepare("insert into(description, value, deposit, wallet_id) values($1, $2, $3, $4) returning id")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	err = stmt.QueryRow(movement.Description, movement.Value, movement.Deposit, movement.Wallet_id).Scan(&movement.Id)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

/*
func GetMovements(wallet_id int64) {
}
*/
