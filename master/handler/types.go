package handler

// A JsonError is the default error message that is generated.
//swagger:response jsonError
type JsonError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// A JsonSuccess is the default success message that is generated.
//swagger:response jsonSuccess
type JsonSuccess struct {
	Message string `json:"message"`
}
