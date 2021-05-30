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

func createWallet(t *testing.T, token string) int64 {
	require := require.New(t)

	walletBody := `{"name": "Fundo thinkpad", "description": "Fundo para comprar novo thinkpad"}`
	req, err := http.NewRequest(http.MethodPost, "/wallet", strings.NewReader(walletBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)

	response := executeRequest(req)
	require.Equal(http.StatusCreated, response.Code, "Should return 201")
	var walletResponse dto.SaveWalletResponseDto
	json.Unmarshal(response.Body.Bytes(), &walletResponse)
	require.NotNil(walletResponse, "Should parse json")

	return walletResponse.Wallet.Id
}

func TestSaveMovement(t *testing.T) {
	clearTables()
	token, err := createAndSignInUser(t)
	require := require.New(t)
	require.Nil(err, "Should return token")

	wallet_id := createWallet(t, token)
	movementBody := `{"description": "Primeiro deposito", "value": 45.0, "deposit": true}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movement", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)

	response := executeRequest(req)
	require.Equal(http.StatusCreated, response.Code, "Should return 201")
	var movementDto dto.AddMovementResponseDto
	json.Unmarshal(response.Body.Bytes(), &movementDto)
	require.NotNil(movementDto.Movement)
	require.Equal(movementDto.Movement.Value, 45.0, "Value should be 45")
}

func TestErrorSaveMovement(t *testing.T) {
        clearTables()
	token, err := createAndSignInUser(t)
	require := require.New(t)
	require.Nil(err, "Should return token")

        wallet_id := createWallet(t, token)
	movementBody := `{"description": "Primeiro deposito", "value": 45.0, "deposit": false}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movement", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)

	response := executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code, "Should return 400")

        movementBody = `{"description": "Primeiro deposito", "value": 20.0, "deposit": null}`
        req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movement", wallet_id), strings.NewReader(movementBody))
        require.Nil(err, "Should create request")
        req.Header.Add("Authorization", "Bearer "+token)
        response = executeRequest(req)
        require.Equal(http.StatusBadRequest, response.Code, "Should return 400")


        movementBody = `{}`
        req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movement", wallet_id), strings.NewReader(movementBody))
        require.Nil(err, "Should create request")
        req.Header.Add("Authorization", "Bearer "+token)
        response = executeRequest(req)
        require.Equal(http.StatusBadRequest, response.Code, "Should return 400")

        movementBody = `{}`
        req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movement", wallet_id), strings.NewReader(movementBody))
        require.Nil(err, "Should create request")
        response = executeRequest(req)
        require.Equal(http.StatusUnauthorized, response.Code, "Should return 401")
}
