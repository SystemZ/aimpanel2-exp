package handler

type JsonError struct {
	ErrorCode int    `json:"error_code" example:"1"`
	Message   string `json:"message" example:""`
}

type JsonSuccess struct {
	Message string `json:"message" example:""`
}
