package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func addClient(t *testing.T) *httptest.ResponseRecorder {
	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
	}
	response := executeRequest(req)

	return response
}

func TestSaveClient(t *testing.T) {
	response := addClient(t)

	if status := response.Code; status != http.StatusCreated {
		t.Errorf("Error saving client, expect status 201 got %v\n", status)
		t.Errorf("%v\n", response.Body.String())
	}
}

func TestErrorSaveClient(t *testing.T) {
	clearTables()
	body := `{"nam": "Fitz", "email": "fitz@gmail.com", "password": "123456"}`
	req, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request %v\n", err)
	}
	response := executeRequest(req)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Error saving client, expect status 400 got %v\n", response.Code)
	}
	addClient(t)
	response = addClient(t)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Error saving client, expect status 400 got %v\n", response.Code)
	}

}
