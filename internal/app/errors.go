package app

import (
	"errors"
	"fmt"
)

type AppError struct {
	Location string
	Code     AppCode
	Err      error
}

func WrapLoc(location string, err *AppError) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Location: fmt.Sprintf("%s: %s", location, err.Location),
		Code:     err.Code,
		Err:      err.Err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Location, e.Err)
	}
	return e.Location
}

func (e *AppError) ErrorOrNil() *string {
	if e == nil {
		return nil
	}

	s := e.Error()
	if s == "" {
		return nil
	}

	return &s
}

func (e *AppError) Is(target error) bool {
	var te *AppError
	if errors.As(target, &te) {
		return e.Code == te.Code
	}
	return false
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) UnwrapAll() error {
	if appErr, ok := e.Err.(*AppError); ok {
		return appErr
	}

	return e.Err
}

func NewAppError(location string, code AppCode, err error) *AppError {
	return &AppError{
		Location: location,
		Code:     code,
		Err:      err,
	}
}

func NewUnknownError(location string, err error) *AppError {
	return &AppError{
		Location: location,
		Code:     CodeUnknown,
		Err:      err,
	}
}

func NewDatabaseError(location string, err error) *AppError {
	return &AppError{
		Location: location,
		Code:     CodeDatabaseError,
		Err:      err,
	}
}
