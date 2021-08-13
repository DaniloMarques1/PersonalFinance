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

func TestSaveMovement(t *testing.T) {
	clearTables()
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")
	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotNil(session, "Session should not be nil")

	token := session.Token

	walletResponse, err := addWallet(token)
	require.Nil(err, "Should have created wallet")
	wallet_id := walletResponse.Wallet.Id

	movementBody := `{"description": "Primeiro deposito", "value": 45.0, "deposit": true}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", wallet_id), strings.NewReader(movementBody))
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
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotNil(session, "Session should not be nil")

	token := session.Token

	walletResponse, err := addWallet(token)
	require.Nil(err, "Should have created wallet")
	wallet_id := walletResponse.Wallet.Id

	movementBody := `{"description": "Primeiro deposito", "value": 45.0, "deposit": false}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)

	response := executeRequest(req)
	require.Equal(http.StatusUnauthorized, response.Code, "Should return 400")

	movementBody = `{"description": "Primeiro deposito", "value": -20.0, "deposit": null}`
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)

	response = executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code, "Should return 400")

	movementBody = `{}`
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)
	response = executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code, "Should return 400")

	movementBody = `{}`
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", wallet_id), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	response = executeRequest(req)
	require.Equal(http.StatusUnauthorized, response.Code, "Should return 401")

	movementBody = `{"description": "Primeiro deposito", "value": 45.0, "deposit": true}`
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("/wallet/%v/movements", -2), strings.NewReader(movementBody))
	require.Nil(err, "Should create request")
	req.Header.Add("Authorization", "Bearer "+token)
	response = executeRequest(req)
	fmt.Println(response.Body.String())
	require.Equal(http.StatusNotFound, response.Code, "Should return 404")
}
