package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danilomarques1/personalfinance/api/dto"
)

type Error struct {
	message string
}

func RespondJson(w http.ResponseWriter, statusCode int, body interface{}) {
	fmt.Println(statusCode)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&body)
}

func HandleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *ApiError:
		apiError := err.(*ApiError)
		RespondJson(w, apiError.Code, dto.ErrorResponseDto{Message: apiError.Message})
	default:
		RespondJson(w, http.StatusInternalServerError, dto.ErrorResponseDto{Message: "Unexpected error"})
	}
}
