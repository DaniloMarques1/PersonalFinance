package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/danilomarques1/personalfinance/api/model"
)

func createAndSignInUser() (string, error) {
	clearTables()
	user := []byte(`{"name": "fitz", "email":"fitz@gmail.com", "password": "123456"}`)
	request, err := http.NewRequest(http.MethodPost, "/client", bytes.NewBuffer(user))
	if err != nil {
		return "", err
	}
	response := executeRequest(request)
	if response.Code != http.StatusCreated {
		return "", fmt.Errorf("Error adding a new user")
	}

	login := []byte(`{"email": "fitz@gmail.com", "password": "123456"}`)
	request, err = http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(login))
	if err != nil {
		return "", err
	}

	response = executeRequest(request)
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

func TestSaveWallet(t *testing.T) {
	token, err := createAndSignInUser()
	if err != nil {
		t.Errorf("Erro creating and signing user %v", err)
	}
	walletBody := []byte(`{"name": "Testando 1", "description": "Description 1"}`)
	request, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(walletBody))
	if err != nil {
		t.Errorf("Error creting wallet request %v", err)
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	if response.Code != http.StatusCreated {
		t.Errorf("Wrong status code, expect 201 got %v", response.Code)
	}

	var wallet model.Wallet
	err = json.Unmarshal(response.Body.Bytes(), &wallet)
	if err != nil {
		t.Errorf("Error unmarshal wallet %v", err)
	}

	if wallet.Name != "Testando 1" {
		t.Errorf("Wrong wallet returned, expect name \"Testando 1\" got %v", wallet.Name)
	}
}
