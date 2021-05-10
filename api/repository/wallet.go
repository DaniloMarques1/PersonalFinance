package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/danilomarques1/personalfinance/api/model"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (wr *WalletRepository) SaveWallet(wallet *model.Wallet) error {
	stmt, err := wr.db.Prepare("insert into wallet(name, description, client_id) values($1, $2, $3) returning id")
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(wallet.Name, wallet.Description, wallet.Client_id).Scan(&wallet.Id)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	return nil
}

func (wr *WalletRepository) RemoveWallet(client_id, wallet_id int64) error {
	tx, err := wr.db.Begin()
	if err != nil {
		log.Printf("Error opening transaction %v\n", err)
		return err
	}

	// cleaning the movements of this specific wallet
	stmt, err := tx.Prepare("delete from movement where wallet_id = $1")
	if err != nil {
		log.Printf("Error preparing query %v\n", err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(wallet_id)
	if err != nil {
		log.Printf("Error removing movements %v\n", err)
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare("delete from wallet where id = $1")
	if err != nil {
		log.Printf("Error preparing query %v\n", err)
		tx.Rollback()
		return fmt.Errorf("Internal server error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(wallet_id)
	if err != nil {
		log.Printf("erro removing wallet %v\n", err)
		tx.Rollback()
		return fmt.Errorf("Internal server error")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		tx.Rollback()
		return fmt.Errorf("Wallet with this id does not exist")
	}

	tx.Commit()

	return nil
}

func (wr *WalletRepository) GetWallets(client_id int64) ([]model.Wallet, error) {
	stmt, err := wr.db.Prepare(`select w.id, w.name, w.description, w.client_id,
                                (select case when sum(value) >= 0 then sum(value) else 0 end from movement where wallet_id=w.id and deposit=true) -
                                (select case when sum(value) >= 0 then sum(value) else 0 end from movement where wallet_id=w.id and deposit=false)
                                as total
                                from wallet as w
                                where client_id=$1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(client_id)
	if err != nil {
		log.Printf("error querying the rows %v\n", err)
		return nil, err
	}

	wallets := make([]model.Wallet, 0)
	for rows.Next() {
		var wallet model.Wallet
		err = rows.Scan(&wallet.Id, &wallet.Name, &wallet.Description, &wallet.Client_id, &wallet.Total)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

/*
select w.name, w.description, c.id, c.name, c.email, c.total, m.description, m.value, m.movement_date,
(select sum(m.value) from tb_movement as m where deposit=true) -
(select case when sum(m.value) >= 0 then sum(m.value) else 0 end from tb_movement as m where deposit=false)
as wallet_total
from tb_wallet as w
join tb_client as c on w.client_id = c.id
join tb_movement as m on m.wallet_id = w.id
where w.id = wallet_id;

// will return the wallet, its owner (client) and the full wallet history (TODO maybe return partially the history)
func (wr *WalletRepository) GetWallet(wallet_id, client_id int64) (*model.Client, *model.Wallet, error) {
	stmt, err := wr.db.Prepare(`select w.id, w.name, w.description, w.total, c.id, c.name, c.email, c.total
                                from tb_wallet as w
                                join tb_client as c on w.client_id = c.id
                                where w.id = $1 and w.client_id = $2`)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil, err
	}
	defer stmt.Close()

	var client model.Client
	var wallet model.Wallet
	err = stmt.QueryRow(wallet_id, client_id).Scan(&wallet.Id, &wallet.Name, &wallet.Description, &wallet.Total,
		&client.Id, &client.Name, &client.Email, &client.Total)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil, err
	}

	return &client, &wallet, nil
}

func (wr *WalletRepository) GetWallets(client_id int64) (*model.Client, []model.Wallet, error) {
	stmt, err := wr.db.Prepare(`select w.id, w.name, w.description, w.total, c.id, c.name, c.email, c.total
                                from tb_wallet as w
                                join tb_client as c
                                on w.client_id = c.id
                                where w.client_id = $1`)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil, err
	}
	defer stmt.Close()

	wallets := make([]model.Wallet, 0)
	rows, err := stmt.Query(client_id)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil, err
	}

	var client model.Client
	for rows.Next() {
		var wallet model.Wallet
		err = rows.Scan(&wallet.Id, &wallet.Name, &wallet.Description, &wallet.Total,
			&client.Id, &client.Name, &client.Email, &client.Total)
        if err != nil {
            return nil, nil, err
        }

		wallets = append(wallets, wallet)
	}

	return &client, wallets, nil
}
*/
