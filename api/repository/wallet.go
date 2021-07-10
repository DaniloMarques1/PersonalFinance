package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/model"
	"github.com/danilomarques1/personalfinance/api/util"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

// save a wallet associating with a client
func (wr *WalletRepository) SaveWallet(wallet *model.Wallet) error {
	stmt, err := wr.db.Prepare("insert into wallet(name, description, client_id) values($1, $2, $3) returning id, created_date")
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(wallet.Name, wallet.Description, wallet.Client_id).Scan(&wallet.Id, &wallet.Created_date)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	return nil
}

// will remove a client wallet and its associated movements
func (wr *WalletRepository) RemoveWallet(wallet_id, client_id int64) error {
	tx, err := wr.db.Begin()
	if err != nil {
		log.Printf("Error opening transaction %v\n", err)
		return err
	}

	// cleaning the movements of this specific wallet
	stmt, err := tx.Prepare("DELETE FROM movement WHERE wallet_id = $1")
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

	stmt, err = tx.Prepare("DELETE FROM wallet WHERE id = $1 AND client_id = $2")
	if err != nil {
		log.Printf("Error preparing query %v\n", err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(wallet_id, client_id)
	if err != nil {
		log.Printf("error removing wallet %v\n", err)
		tx.Rollback()
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		tx.Rollback()
                log.Printf("Rows affected %v", rowsAffected)
		return util.NewApiError("Wallet not found", http.StatusNotFound)
	}

	tx.Commit()

	return nil
}

// will return all wallets of a client and a total which is the summ of all
// movements (deposit) of a wallet
func (wr *WalletRepository) FindAll(client_id int64) ([]model.Wallet, float64, error) {
	stmt, err := wr.db.Prepare(`select w.id, w.name, w.description, w.created_date, w.client_id,
                                (select case when sum(value) >= 0 then sum(value) else 0 end from movement where wallet_id=w.id and deposit=true) -
                                (select case when sum(value) >= 0 then sum(value) else 0 end from movement where wallet_id=w.id and deposit=false)
                                as total
                                from wallet as w
                                where client_id=$1`)
	if err != nil {
                log.Printf("Error preparing findAll query %v", err)
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(client_id)
	if err != nil {
		log.Printf("error querying the rows %v\n", err)
		return nil, 0, err
	}
	defer rows.Close()

	wallets := make([]model.Wallet, 0)
        var total float64
	for rows.Next() {
		var wallet model.Wallet
		err = rows.Scan(&wallet.Id, &wallet.Name, &wallet.Description, &wallet.Created_date, &wallet.Client_id, &wallet.Total)
		if err != nil {
                        log.Printf("Error scanning rows %v", err)
			return nil, 0, err
		}
		wallets = append(wallets, wallet)
		total += wallet.Total
	}

	return wallets, total, nil
}

func (wr *WalletRepository) FindById(wallet_id, client_id int64) (*model.Wallet, error) {
	stmt, err := wr.db.Prepare("SELECT id FROM wallet WHERE id = $1 AND client_id = $2)")
	if err != nil {
                log.Printf("Error preparing find by id %v", err)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(wallet_id, client_id)
	var wallet model.Wallet
	err = row.Scan(&wallet.Id)
	if err != nil {
                log.Printf("Error scanning row %v", err)
		return nil, err
	}

	return &wallet, nil
}
