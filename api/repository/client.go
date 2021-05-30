package repository

import (
	"database/sql"
	"log"

	"github.com/danilomarques1/personalfinance/api/model"
)

type ClientRepository struct {
	Db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		Db: db,
	}
}

func (cr *ClientRepository) SaveClient(client *model.Client) error {
	stmt, err := cr.Db.Prepare("insert into client(name, email, password_hash) values($1, $2, $3) returning id")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(&client.Name, &client.Email, &client.PasswordHash).Scan(&client.Id)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func (cr *ClientRepository) FindById(id int64) (*model.Client, error) {
	stmt, err := cr.Db.Prepare("select * from client where id = $1")
	if err != nil {
		log.Printf("Error prepating statement %v", err)
		return nil, err
	}
	defer stmt.Close()

	var client model.Client
	err = stmt.QueryRow(id).Scan(&client.Id, &client.Name, &client.Email, &client.PasswordHash)
	if err != nil {
		log.Printf("Error finding by id %v", err)
		return nil, err
	}

	return &client, nil
}

func (cr *ClientRepository) FindByEmail(email string) (*model.Client, error) {
	stmt, err := cr.Db.Prepare("select * from client where email = $1")
	if err != nil {
		log.Printf("Error preparing statement %v", err)
		return nil, err
	}
	defer stmt.Close()

	var client model.Client
	err = stmt.QueryRow(email).Scan(&client.Id, &client.Name, &client.Email, &client.PasswordHash)
	if err != nil {
		log.Printf("Error finding email %v", err)
		return nil, err
	}

	return &client, nil
}
