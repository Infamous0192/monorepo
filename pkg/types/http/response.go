package http

import "app/pkg/types/pagination"

// GeneralResponse represents a standard API response
type GeneralResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Metadata pagination.Metadata `json:"metadata"`
	Result   []interface{}       `json:"result"`
}
