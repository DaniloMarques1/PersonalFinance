package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danilomarques1/personalfinance/api/server"
	_ "github.com/lib/pq"
)

var App server.App

func TestMain(m *testing.M) {
	App.Initialize(server.DbConn{
		Host:     "0.0.0.0",
		Port:     "5433",
		User:     "fitz",
		Password: "123456",
		DbName:   "personalwallet",
	})

	code := m.Run()
	clearTables()
	os.Exit(code)
}

func clearTables() {
	tx, err := App.Db.Begin()
	if err != nil {
		log.Fatalf("Error creating the transaction %v\n", err)
	}
	_, err = tx.Exec("truncate table movement cascade;")
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("truncate table wallet cascade;")
	if err != nil {
		tx.Rollback()
	}

	_, err = tx.Exec("truncate table client cascade;")
	if err != nil {
		tx.Rollback()
	}

	tx.Commit()
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, request)

	return rr
}
