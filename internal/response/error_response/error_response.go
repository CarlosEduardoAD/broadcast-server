package error_response

import (
	"errors"
	"time"
)

type ErrorResponse struct {
	Error     error
	CreatedAt time.Time
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error:     err,
		CreatedAt: time.Now(),
	}
}

func (er *ErrorResponse) Validate() error {
	if er.Error == nil {
		return errors.New("err should not be empty")
	}

	return nil
}
