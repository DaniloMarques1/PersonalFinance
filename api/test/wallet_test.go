package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danilomarques1/personalfinance/api/dto"
)

func createAndSignInUser() (string, error) {
	clearTables()
	user := `{"name": "fitz", "email":"fitz@gmail.com", "password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(user))
	if err != nil {
		return "", err
	}
	response := executeRequest(request)
	if response.Code != http.StatusCreated {
		return "", fmt.Errorf("Error adding a new user")
	}

	login := `{"email": "fitz@gmail.com", "password": "123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(login))
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
		t.Fatalf("Erro creating and signing user %v", err)
	}
	walletBody := `{"name": "Testando 1", "description": "Description 1"}`
	request, err := http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(walletBody))
	if err != nil {
		t.Fatalf("Error creting wallet request %v", err)
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	if response.Code != http.StatusCreated {
		t.Fatalf("Wrong status code, expect 201 got %v", response.Code)
	}

	var walletDto dto.SaveWalletResponseDto
	err = json.Unmarshal(response.Body.Bytes(), &walletDto)
	if err != nil {
		t.Fatalf("Error unmarshal wallet %v", err)
	}

	wallet := walletDto.Wallet
	if wallet.Name != "Testando 1" {
		t.Fatalf("Wrong wallet returned, expect name \"Testando 1\" got %v", wallet.Name)
	}
}
