package util

type ApiError struct {
	Message string
	Code    int
}

func NewApiError(message string, code int) *ApiError {
	return &ApiError{
		Message: message,
		Code:    code,
	}
}

func (apiError *ApiError) Error() string {
	return apiError.Message
}

/*
type EmailAlreadyTaken struct {
        ApiError
}

func NewEmailAlreadyTaken(message string, code int) *EmailAlreadyTaken {
        return &EmailAlreadyTaken{
                Message: message,
                Code: code,
        }
}
*/
