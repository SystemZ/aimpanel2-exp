package response

// A JsonError is the default error message that is generated.
//swagger:response jsonError
type JsonError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type JsonSuccess struct {
	Message string `json:"message"`
}
