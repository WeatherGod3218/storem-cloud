package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
