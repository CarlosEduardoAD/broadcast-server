package utils

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBypassCheckAlwaysTrue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	ok := BypassCheck(req)

	assert.True(t, ok, "expected BypassCheck to always return true")
}

func TestGenerateRandomIpDeterministic(t *testing.T) {
	// reproduce the same random source used in the implementation
	expectedRand := rand.New(rand.NewSource(99))
	expected := "192.168.0." + fmt.Sprint(expectedRand.Int63())

	got := GenerateRandomIp()

	assert.Equal(t, expected, got, "expected GenerateRandomIp to return deterministic value with fixed seed")
}
