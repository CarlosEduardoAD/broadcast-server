package error_response

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenGivenAInvalidErrorResponseItShouldReturnError(t *testing.T) {
	errResp := NewErrorResponse(nil)

	err := errResp.Validate()
	assert.NotNil(t, err)
}

func TestWhenGivenAValidErrorResponseItShouldSucceed(t *testing.T) {
	errResp := NewErrorResponse(errors.New("sample error"))

	err := errResp.Validate()
	assert.Nil(t, err)
}
