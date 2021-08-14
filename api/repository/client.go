package repository

import (
	"database/sql"
	"log"

	"github.com/danilomarques1/personalfinance/api/model"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (cr *ClientRepository) SaveClient(client *model.Client) error {
	stmt, err := cr.db.Prepare("insert into client(name, email, password_hash) values($1, $2, $3) returning id")
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
	stmt, err := cr.db.Prepare("select * from client where id = $1")
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
	stmt, err := cr.db.Prepare("select * from client where email = $1")
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

func (cr *ClientRepository) UpdateClient(client *model.Client) error {
	stmt, err := cr.db.Prepare("update client set name = $1, email = $2, password_hash = $3 where id = $4")
	if err != nil {
		log.Printf("Error preparing update %v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(client.Name, client.Email, client.PasswordHash, client.Id)
	if err != nil {
		log.Printf("Error executing update %v\n", err)
		return err
	}

	return nil
}
