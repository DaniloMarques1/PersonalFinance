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
	addClient(t, "Fitz", "fitz@gmail.com", "123456")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(t, err, "Err should be nil")
	require.NotNil(t, session, "Session should not be nil")

	token := session.Token
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
	addClient(t, "Fitz", "fitz@gmail.com", "123456")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(t, err, "Err should be nil")
	require.NotNil(t, session, "Session should not be nil")
	token := session.Token

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
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotNil(session, "Session should not be nil")

	token := session.Token

	walletResponse, err := addWallet(token)
	require.Nil(err, "Error should be nil")
	require.NotNil(walletResponse, "Id should not be -1")
	wallet_id := walletResponse.Wallet.Id

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/wallet/%v", wallet_id), nil)
	require.Nil(err, "Error creating request")
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)

	require.Equal(http.StatusNoContent, response.Code)
}

func TestUpdateWallet(t *testing.T) {
	clearTables()
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")
	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotEqual(session.Token, "")

	token := session.Token
	walletResponse, err := addWallet(token)
	require.Nil(err, "Should have created wallet")

	updateBody := `{"name": "New wallet name", "description": "new wallet description"}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/wallet/%v", walletResponse.Wallet.Id), strings.NewReader(updateBody))
	req.Header.Add("Authorization", "Bearer "+token)
	require.Nil(err, "Should have created the request")

	response := executeRequest(req)
	fmt.Println(response.Body.String())
	require.Equal(http.StatusNoContent, response.Code)
}

func TestErrorUpdateWallet(t *testing.T) {
	clearTables()
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")
	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotEqual(session.Token, "")

	token := session.Token
	walletResponse, err := addWallet(token)
	require.Nil(err, "Should have created wallet")

	updateBody := `{"ame": "New wallet name", "description": "new wallet description"}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/wallet/%v", walletResponse.Wallet.Id), strings.NewReader(updateBody))
	req.Header.Add("Authorization", "Bearer "+token)
	require.Nil(err, "Should have created the request")
	response := executeRequest(req)

	require.Equal(http.StatusBadRequest, response.Code)

	updateBody = `{"name": "New wallet name", "description": "new wallet description"}`
	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("/wallet/%v", walletResponse.Wallet.Id), strings.NewReader(updateBody))
	require.Nil(err, "Should have created the request")
	response = executeRequest(req)

	require.Equal(http.StatusUnauthorized, response.Code)
}
