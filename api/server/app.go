package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/handler"
	"github.com/danilomarques1/personalfinance/api/repository"
	"github.com/danilomarques1/personalfinance/api/util"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

type DbConn struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

const tables = `
    CREATE TABLE IF NOT EXISTS client (
        id SERIAL PRIMARY KEY NOT NULL,
        name VARCHAR(60) NOT NULL,
        email VARCHAR(60) NOT NULL UNIQUE,
        password_hash VARCHAR(100) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS wallet (
        id SERIAL PRIMARY KEY NOT NULL,
        name VARCHAR(40) NOT NULL,
        description VARCHAR(120) NOT NULL,
        client_id INT NOT NULL,
        created_date TIMESTAMP DEFAULT NOW(),
        CONSTRAINT fk_client_id FOREIGN KEY(client_id) REFERENCES client(id)
    );

    CREATE TABLE IF NOT EXISTS movement(
        id SERIAL NOT NULL PRIMARY KEY,
        description VARCHAR(150),
        value DECIMAL NOT NULL,
        deposit BOOLEAN NOT NULL DEFAULT TRUE,
        movement_date timestamp not null default now(),
        wallet_id INT NOT NULL,
        CONSTRAINT fk_wallet_id FOREIGN KEY(wallet_id) REFERENCES wallet(id)
    );
`

// returns a postgres connection string based on the values
// of the DbConn passed
func connString(db DbConn) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DbName)
}

// set up the web servers database connection
// and the servers routes
func (app *App) Initialize(db DbConn) {
	con := connString(db)
	var err error
	app.Db, err = sql.Open("postgres", con)
	if err != nil {
		log.Fatalf("Error opening database conection %v", err)
	}
	if _, err := app.Db.Exec(tables); err != nil {
		log.Fatalf("Error creating tables %v", err)
	}
	app.Router = mux.NewRouter()

	clientRepository := repository.NewClientRepository(app.Db)
	clientHandler := handler.NewClientHandler(clientRepository)

	walletRepository := repository.NewWalletRepository(app.Db)
	walletHandler := handler.NewWalletHandler(walletRepository)

	//movementRepository := repository.NewMovementRepository(app.Db)
	//movementHandler := handler.NewMovementHandler(movementRepository)

	app.Router.HandleFunc("/client", clientHandler.SaveClient).Methods(http.MethodPost)
	app.Router.HandleFunc("/session", clientHandler.CreateSession).Methods(http.MethodPost)

	app.Router.Handle("/wallet", util.AuthorizationMiddleware(http.HandlerFunc(walletHandler.SaveWallet))).Methods(http.MethodPost)
	app.Router.Handle("/wallet/{wallet_id}", util.AuthorizationMiddleware(http.HandlerFunc(walletHandler.RemoveWallet))).Methods(http.MethodDelete)

	app.Router.Handle("/wallet/", util.AuthorizationMiddleware(http.HandlerFunc(walletHandler.GetWallets))).Methods(http.MethodGet)

	// TODO add route to retrieve movements of a specific wallet

	// TODO add routes to deposit
	// TODO add routes to withdraw
}

// starts the web server
func (app *App) Start() {
	server := &http.Server{
		Handler: app.Router,
		Addr:    "127.0.0.1:8080",
	}

	log.Fatal(server.ListenAndServe())
}
