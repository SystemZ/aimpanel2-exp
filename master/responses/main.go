package responses

type JsonError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
