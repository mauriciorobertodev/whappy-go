package http

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app"
)

type HttpResponse struct {
	Code    app.AppCode `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
	Errors  *ErrorBag   `json:"errors,omitempty"`
}

func NewResponse(code app.AppCode, message string, data interface{}, err *app.AppError) *HttpResponse {
	return &HttpResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   err.ErrorOrNil(),
	}
}

func NewErrorResponse(message string, err *app.AppError) *HttpResponse {
	return &HttpResponse{
		Code:    err.Code,
		Message: message,
		Data:    nil,
		Error:   err.ErrorOrNil(),
		Errors:  nil,
	}
}

func NewInternalErrorResponse(location string, err error) *HttpResponse {
	appErr := app.NewAppError(location, app.CodeUnknown, err)
	return &HttpResponse{
		Code:    appErr.Code,
		Message: "Internal server error",
		Data:    nil,
		Error:   appErr.ErrorOrNil(),
		Errors:  nil,
	}
}

func NewSuccessResponse(message string, data interface{}) *HttpResponse {
	return &HttpResponse{
		Code:    app.CodeSuccess,
		Message: message,
		Data:    data,
		Error:   nil,
		Errors:  nil,
	}
}

func NewSuccessEmptyResponse() *HttpResponse {
	return &HttpResponse{
		Code:    app.CodeSuccess,
		Message: "Success",
		Data:    nil,
		Error:   nil,
		Errors:  nil,
	}
}

func NewValidationErrorResponse(errors *ErrorBag) *HttpResponse {
	return &HttpResponse{
		Code:    app.CodeValidationFailed,
		Message: "The given data was invalid.",
		Data:    nil,
		Error:   nil,
		Errors:  errors,
	}
}

func NewInvalidJSONResponse() *HttpResponse {
	appError := app.NewAppError(
		"invalid JSON format",
		app.CodeInvalidJSON,
		nil,
	)

	return NewErrorResponse("Invalid JSON format", appError)
}
