package core

import "fmt"

type ApiError struct {
	Opcode           int                    `json:"opcode"`
	ErrorStr         string                 `json:"error"`
	MessageStr       string                 `json:"message"`
	LocalizedMessage string                 `json:"localizedMessage"`
	Title            string                 `json:"title"`
	Payload          map[string]interface{} `json:"payload"`
}

func (e *ApiError) Error() string {
	msg := e.MessageStr
	if msg == "" {
		msg = e.ErrorStr
	}
	return fmt.Sprintf("ApiError(opcode=%d): %s", e.Opcode, msg)
}

func (e *ApiError) ErrorCode() string {
	return e.ErrorStr
}

func NewApiError(opcode int, errStr, message, localizedMessage, title string, payload map[string]interface{}) *ApiError {
	return &ApiError{
		Opcode:           opcode,
		ErrorStr:         errStr,
		MessageStr:       message,
		LocalizedMessage: localizedMessage,
		Title:            title,
		Payload:          payload,
	}
}

type PyMaxError struct {
	Message string
}

func (e *PyMaxError) Error() string {
	return fmt.Sprintf("PyMaxError: %s", e.Message)
}

func NewPyMaxError(msg string) *PyMaxError {
	return &PyMaxError{Message: msg}
}

type UploadError struct {
	Message string
}

func (e *UploadError) Error() string {
	return fmt.Sprintf("UploadError: %s", e.Message)
}

func NewUploadError(msg string) *UploadError {
	return &UploadError{Message: msg}
}
