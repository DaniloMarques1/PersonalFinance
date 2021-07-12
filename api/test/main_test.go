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
		Port:     "5432",
		User:     "fitz",
		Password: "123456",
		DbName:   "personalwallet",
	})

	code := m.Run()
	clearTables()
	os.Exit(code)
}

// will delete all data in the database
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

// will add a new client and return its dto
func addClient(t *testing.T, name, email, password string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"name": "%v", "email": "%v", "password": "%v"}`, name, email, password)
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
		t.FailNow()
	}
	response := executeRequest(req)

	return response
}

// execute a request
func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, request)

	return rr
}

// returns a token
func signIn(email, password string) (string, error) {
	loginBody := fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password)
	request, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(loginBody))
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

// add a new wallet
func addWallet(token string) (*dto.SaveWalletResponseDto, error) {
	body := `{"name":"Lista do natal", "description": "Carteira onde será o salvo dinheiro do natal"}`
	request, err := http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)

	response := executeRequest(request)
	if response.Code != http.StatusCreated {
		return nil, fmt.Errorf("Wrong status returned")
	}

	var walletResponse dto.SaveWalletResponseDto
	err = json.Unmarshal(response.Body.Bytes(), &walletResponse)
	if err != nil {
		return nil, err
	}

	return &walletResponse, nil
}
