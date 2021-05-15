package dto

type AddMovementDto struct {
	Description string  `json:"description"`
	Deposit     bool    `json:"deposit"`
	Value       float64 `json:"value"`
}
