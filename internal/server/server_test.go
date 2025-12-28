package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateServerWithValidData(t *testing.T) {
	s := NewServer(12)

	assert.NotNil(t, s, "expected server instance to be created")
}
