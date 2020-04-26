package response

type Token struct {
	Token string `json:"token"`
}

type JsonError struct {
	ErrorCode int    `json:"error_code" example:"1"`
	Message   string `json:"message" example:""`
}
