package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/stretchr/testify/require"
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

        require.Nil(t, err, "Should create a user")

	walletBody := `{"name": "Testando 1", "description": "Description 1"}`
	request, err := http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(walletBody))
        require.Nil(t, err, "Error creating wallet request")

	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)

        require.Equal(t, response.Code, http.StatusCreated, "Should return 201 when adding a wallet")

	var walletDto dto.SaveWalletResponseDto
	err = json.Unmarshal(response.Body.Bytes(), &walletDto)
        require.Nil(t, err, "Should parse response")

	wallet := walletDto.Wallet
        require.Equal(t, wallet.Name, "Testando 1", "Should return the correct wallet name")
}

func TestErrorSaveWallet(t *testing.T) {
	token, err := createAndSignInUser()
        require.Nil(t, err, "Should create a user")

	walletBody := `{"name": "", "description": "Description 1"}`
	request, err := http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(walletBody))
        require.Nil(t, err, "Error creating wallet request")

	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
        require.Equal(t, http.StatusBadRequest, response.Code, "Should return http 400")

        walletBody = `{"name": "Testando 1", "description": "Description 1"}`
        request, err = http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(walletBody))
        require.Nil(t, err, "Should create request")
        response = executeRequest(request)
        require.Equal(t, http.StatusUnauthorized, response.Code)
        
}
