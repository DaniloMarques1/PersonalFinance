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

func TestSaveWallet(t *testing.T) {
	token, err := createAndSignInUser(t)

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
	token, err := createAndSignInUser(t)
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

func TestRemoveWallet(t *testing.T) {
	clearTables()
	token, err := createAndSignInUser(t)
	require.Nil(t, err, "Should have created and signed the user")
	require.NotEqual(t, token, "", "Should have returned a token")

	wallet_id, err := addWallet(token)
	require.Nil(t, err, "Error should be nil")
	require.NotEqual(t, -1, wallet_id, "Id should not be -1")

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/wallet/%v", wallet_id), nil)
	require.Nil(t, err, "Error creating request")
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	require.Equal(t, http.StatusNoContent, response.Code)
}
