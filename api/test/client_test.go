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
	response := addClient(t, "Fitz", "fitz@gmail.com", "123456")

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

	addClient(t, "Fitz", "fitz@gmail.com", "123456")
	response = addClient(t, "Fitz", "fitz@gmail.com", "123456")

	require.Equal(http.StatusBadRequest, response.Code, "Should not let register the client using the same email")

	body = `{"name": "Danilo Marques de oliveira Danilo Marques de Oliveira Danilo Marques de Oliveira"
		, "email": "danilo@hotmail.com", "password": "danilodanilodanilodanilodanilo"}`
	req, err = http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	require.Nil(err, "Err should be nil")

	response = executeRequest(req)
	require.Equal(http.StatusBadRequest, response.Code,
		"Should not let register an user with more than 60 chars in his name/email or over 20 in his password")
}

func TestCreateSession(t *testing.T) {
	clearTables()
	addClient(t, "Fitz", "fitz@gmail.com", "123456")

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
	require.NotNil(sessionResponse.RefreshToken, "When creating session should return refresh token")
	require.NotNil(sessionResponse.Client, "When creating a session should return the client")

	require.Equal("fitz@gmail.com", sessionResponse.Client.Email, "Email shoud match")
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

func TestRefreshToken(t *testing.T) {
	clearTables()
	require := require.New(t)

	addClient(t, "Fitz", "fitz@gmail.com", "123456")
	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	require.NotNil(session, "Err should be nil")

	req, err := http.NewRequest(http.MethodPut, "/session", nil)
	require.Nil(err, "Err Should be nil")
	req.Header.Add("refresh_token", session.RefreshToken)

	response := executeRequest(req)
	require.Equal(http.StatusOK, response.Code)
	var nSession dto.SessionResponseDto
	json.Unmarshal(response.Body.Bytes(), &nSession)

	require.NotNil(nSession, "Session should be the returned body")
	require.NotNil(nSession.RefreshToken, "should have a refresh token")
	require.NotNil(nSession.Token, "Should have a token")
	require.NotNil(nSession.Client, "Should have a client")

	fmt.Printf("session %v\n", nSession)
}

func TestUpdateClient(t *testing.T) {
	clearTables()
	require := require.New(t)

	response := addClient(t, "Fitz", "fitz@gmail.com", "123456")
	require.Equal(response.Code, http.StatusCreated, "Status should be 201")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	token := session.Token

	body := `{"name": "Fitz Calvary", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`

	request, err := http.NewRequest(http.MethodPut, "/client", strings.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+token)
	require.Nil(err, "Error Should be nil")

	response = executeRequest(request)
	require.Equal(http.StatusNoContent, response.Code, "Status should be 204")
}

func TestErrorUpdateClient(t *testing.T) {
	clearTables()
	require := require.New(t)

	name := "Fitz"
	email := "fitz@gmail.com"
	password := "123456"
	confirmPassword := "different_password"

	response := addClient(t, "Fitz", "fitz@gmail.com", "123456")
	require.Equal(response.Code, http.StatusCreated, "Status should be 201")

	session, err := signIn("fitz@gmail.com", "123456")
	require.Nil(err, "Err should be nil")
	token := session.Token

	body := fmt.Sprintf(`{"name": "%v", "email": "%v", "password": "%v", "confirm_password": "%v"}`,
		name, email, password, confirmPassword)
	request, err := http.NewRequest(http.MethodPut, "/client", strings.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+token)
	require.Nil(err, "Error Should be nil")

	response = executeRequest(request)
	require.Equal(http.StatusBadRequest, response.Code, "Should return 400")
}
