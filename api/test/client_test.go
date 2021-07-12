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

func TestSaveClient(t *testing.T) {
	clearTables()
	response := addClient(t)

	require.Equal(t, http.StatusCreated, response.Code, fmt.Sprintf("Error saving client, expect 201 got %v", response.Code))
}

func TestErrorSaveClient(t *testing.T) {
	clearTables()
	require := require.New(t)

	body := `{"nam": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	require.Nil(err, "Err should be nil")

	response := executeRequest(req)

	require.Equal(http.StatusBadRequest, response.Code, "Should not register an user passing an invalid body")

	addClient(t)
	response = addClient(t)

	require.Equal(http.StatusBadRequest, response.Code, "Should not let register the client using the same email")

	body = `{"name": "Danilo Marques de oliveira Danilo Marques de Oliveira Danilo Marques de Oliveira", "email": "danilo@hotmail.com", "password": "danilodanilodanilodanilodanilo"}`
	req, err = http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	require.Nil(err, "Err should be nil")

	response = executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code, "Should not let register an user with more than 60 chars in his name/email or over 20 in his password")
}

func TestCreateSession(t *testing.T) {
	clearTables()
	addClient(t)

	require := require.New(t)
	body := `{"email": "fitz@gmail.com", "password": "123456"}`

	req, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	require.Nil(err, "Err should be nil")
	response := executeRequest(req)

	require.Equal(http.StatusOK, response.Code, "Should return 200 when creating a new session")

	var sessionResponse dto.SessionResponseDto
	json.Unmarshal(response.Body.Bytes(), &sessionResponse)

	require.NotNil(sessionResponse, "sessionResponse must not be nil")
	require.NotNil(sessionResponse.Token, "When creating a session should return a token")
	require.NotNil(sessionResponse.Client, "When creating a session should return the client")

}

func TestErrorCreateSession(t *testing.T) {
	clearTables()
	require := require.New(t)

	body := `{"email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
        require.Nil(err, "Err shoudl be nil")

	response := executeRequest(req)
	require.Equal(http.StatusUnauthorized, response.Code, "Should return a unauthorized status")

	body = `{"email": "", "password": ""}`
	req, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
        require.Nil(err, "Err shoudl be nil")

	response = executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code, "Should return bad request")
}

func TestUpdateClient(t *testing.T) {
	clearTables()
	require := require.New(t)

	response := addClient(t)
	require.Equal(response.Code, http.StatusCreated, "Status should be 201")

	body := `{"name": "Fitz Calvary", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`

	request, err := http.NewRequest(http.MethodPut, "/client", strings.NewReader(body))
	require.Nil(err, "Error Should be nil")

	response = executeRequest(request)
	require.Equal(response.Code, http.StatusNoContent, "Status should be 204")
}
