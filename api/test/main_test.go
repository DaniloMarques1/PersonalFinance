package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/danilomarques1/personalfinance/api/dto"
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

func addClient(t *testing.T) *httptest.ResponseRecorder {
	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
	}
	response := executeRequest(req)

	return response
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, request)

	return rr
}

// returns a token
func createAndSignInUser(t *testing.T) (string, error) {
	clearTables()
	addClient(t)

	login := `{"email": "fitz@gmail.com", "password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(login))
	if err != nil {
		return "", err
	}

	response := executeRequest(request)
	if response.Code != http.StatusOK {
		return "", fmt.Errorf("Error Sign in")
	}
	var session dto.SessionResponseDto
	err = json.Unmarshal(response.Body.Bytes(), &session)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error getting response")
	}

	return session.Token, nil
}
