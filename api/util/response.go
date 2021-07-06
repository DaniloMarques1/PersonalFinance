package util

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	message string
}

func RespondJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&body)
}

// TODO how to handle errors

func HandleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case ApiError:
                RespondJson(w, err.Code, dto.ErrorResponse{Message: err.Message})
	default:
		RespondJson(w, http.StatusInternalServerError, dto.ErrorResponse{Message: "Unexpected error"})
	}
}
