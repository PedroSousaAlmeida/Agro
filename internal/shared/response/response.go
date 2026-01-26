package response

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{Data: data}
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Message: message}
}

func NewErrorResponseWithCode(message, code string) ErrorResponse {
	return ErrorResponse{Message: message, Code: code}
}
