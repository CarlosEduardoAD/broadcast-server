package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenCreatingANewInvalidMessageItShouldReturnError(t *testing.T) {
	msg := NewMessage("")

	err := msg.Validate()
	assert.NotNil(t, err)
}

func TestWhenCreatingANewValidMessageItShouldSucceed(t *testing.T) {
	msg := NewMessage("Hello, World!")

	err := msg.Validate()
	assert.Nil(t, err)
}
