package dto

// GenericResponse represents a standard structure for API responses.
// It includes a status code, a message, and an optional data payload.
//
// Fields:
// - StatusCode: The HTTP status code representing the result of the request (required).
// - Message: A message providing more details about the response (required).
// - Data: Optional field for including additional data in the response. It can be any type and will be omitted if not present.
type GenericResponse struct {
	StatusCode int         `json:"status_code" validate:"required"`
	Message    string      `json:"message" validate:"required"`
	Data       interface{} `json:"data,omitempty"`
}
