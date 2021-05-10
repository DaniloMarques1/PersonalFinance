package test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danilomarques1/personalfinance/api/server"
	_ "github.com/lib/pq"
)

var app server.App

func TestMain(m *testing.M) {
	app.Initialize(server.DbConn{
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
	tx, err := app.Db.Begin()
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
	app.Router.ServeHTTP(rr, request)

	return rr
}

func addClient(t *testing.T) *httptest.ResponseRecorder {
	body := []byte(`{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`)
	req, err := http.NewRequest(http.MethodPost, "/client", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
	}
	response := executeRequest(req)

	return response
}

func TestSaveClient(t *testing.T) {
	response := addClient(t)

	if status := response.Code; status != http.StatusCreated {
		t.Errorf("Error saving client, expect status 201 got %v\n", status)
		t.Errorf("%v\n", response.Body.String())
	}
}

func TestErrorSaveClient(t *testing.T) {
	clearTables()
	body := []byte(`{"nam": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`)
	req, err := http.NewRequest(http.MethodPost, "/client", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
	}
	response := executeRequest(req)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Error saving client, expect status 400 got %v\n", response.Code)
	}
	addClient(t)
	response = addClient(t)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Error saving client, expect status 400 got %v\n", response.Code)
	}

}
