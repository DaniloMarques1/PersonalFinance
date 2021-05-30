package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danilomarques1/personalfinance/api/dto"
	"github.com/stretchr/testify/assert"
)

func TestSaveClient(t *testing.T) {
	clearTables()
	response := addClient(t)

	assert.Equal(t, http.StatusCreated, response.Code, fmt.Sprintf("Error saving client, expect 201 got %v", response.Code))
}

func TestErrorSaveClient(t *testing.T) {
	clearTables()
	body := `{"nam": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request %v\n", err)
	}
	assert := assert.New(t)

	response := executeRequest(req)

	assert.Equal(http.StatusBadRequest, response.Code, "Should not register an user passing an invalid body")

	addClient(t)
	response = addClient(t)

	assert.Equal(http.StatusBadRequest, response.Code, "Should not let register the client using the same email")

	body = `{"name": "Danilo Marques de oliveira Danilo Marques de Oliveira Danilo Marques de Oliveira", "email": "danilo@hotmail.com", "password": "danilodanilodanilodanilodanilo"}`
	req, err = http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request %v", err)
	}

	response = executeRequest(req)
	assert.Equal(http.StatusBadRequest, response.Code, "Should not let register an user with more than 60 chars in his name/email or over 20 in his password")
}

func TestCreateSession(t *testing.T) {
	clearTables()
	addClient(t)

	assert := assert.New(t)
	body := `{"email": "fitz@gmail.com", "password": "123456"}`

	req, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request %v", err)
	}
	response := executeRequest(req)

	assert.Equal(http.StatusOK, response.Code, "Should return 200 when creating a new session")

	var sessionResponse dto.SessionResponseDto
	json.Unmarshal(response.Body.Bytes(), &sessionResponse)

	assert.NotNil(sessionResponse.Token, "When creating a session should return a token")
	assert.NotNil(sessionResponse.Client, "When creating a session should return the client")

}

func TestErrorCreateSession(t *testing.T) {
	clearTables()
	assert := assert.New(t)

	body := `{"email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request %v", err)
	}

	response := executeRequest(req)
	assert.Equal(http.StatusUnauthorized, response.Code, "Should return a unauthorized status")

	body = `{"email": "", "password": ""}`
	req, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Error creating request %v", err)
	}

	response = executeRequest(req)
	assert.Equal(http.StatusBadRequest, response.Code, "Should return bad request")
}
