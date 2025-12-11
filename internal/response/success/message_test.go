package success

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenCreatingANewMessageWithInvalidDataItShouldError(t *testing.T) {
	msg := NewMessage("", "user1")

	err := msg.Validate()

	assert.NotNil(t, err)
}

func TestWhenCreatingANewMessageWithValidDataItShouldSucceed(t *testing.T) {
	msg := NewMessage("Hello, World!", "user1")

	err := msg.Validate()

	assert.Nil(t, err)
}
