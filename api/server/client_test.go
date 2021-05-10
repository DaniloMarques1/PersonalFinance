package server

import (
	//"bytes"
	//"net/http"
	//"net/http/httptest"
	"os"
	"testing"

	"github.com/danilomarques1/personalfinance/api/server"
)

var app server.App

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestSaveClient(t *testing.T) {
	/*
	   body := []byte(`{name: "Fitz", "email": "fitz@gmail.com", "password": "123456"}`)
	   req, err := http.NewRequest(http.MethodPost, "/client", bytes.NewBuffer(body))
	   if err != nil {
	       t.Errorf("Error creating request %v\n", err)
	   }
	   rr := httptest.NewRecorder()
	   handler := http.HandleFunc(SaveClient)
	   handler.ServeHTTP(rr, req)

	   if status := rr.Code; status != http.StatusCreated {
	       t.Errorf("Error saving client, expect status 201 got %v\n", status)
	   }
	*/
}
