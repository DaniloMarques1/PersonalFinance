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
	stmt, err := mr.db.Prepare("insert into movement(description, value, deposit, wallet_id) values($1, $2, $3, $4) returning id, movement_date")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(movement.Description, movement.Value, movement.Deposit, movement.Wallet_id).Scan(&movement.Id, &movement.MovementDate)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func (mr *MovementRepository) CanWithDraw(wallet_id int64, value float64) (bool, error) {
	stmt, err := mr.db.Prepare(`SELECT COALESCE(SUM(value), 0) from movement where deposit = true and wallet_id = $1`)
	if err != nil {
		log.Printf("%v", err)
		return false, err
	}
	defer stmt.Close()

	var total float64
	err = stmt.QueryRow(wallet_id).Scan(&total)
	if err != nil {
		log.Printf("Error scanning %v", err)
		return false, err
	}

	if value <= total {
		return true, nil
	}

	return false, nil

}

func (mr *MovementRepository) FindAll(wallet_id int64) ([]model.Movement, error) {
	stmt, err := mr.db.Prepare(`select id, description, deposit, value, movement_date, wallet_id 
                                    from movement
                                    where wallet_id = $1
                                    order by movement_date desc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(wallet_id)
	if err != nil {
		log.Printf("Error querying movements %v", err)
		return nil, err
	}
	movements := make([]model.Movement, 0)
	for rows.Next() {
		var movement model.Movement
		err = rows.Scan(&movement.Id, &movement.Description, &movement.Deposit, &movement.Value, &movement.MovementDate, &movement.Wallet_id)
		if err != nil {
			log.Printf("Error scanning %v", err)
			return nil, err
		}

		movements = append(movements, movement)
	}

	return movements, nil
}
