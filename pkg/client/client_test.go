package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	spyDialer = &SpyDialer{}
)

func TestWhenGivenAWrongClientItShouldReturnError(t *testing.T) {
	var invalid any = "invalid"

	defer func() {
		if r := recover(); r != nil {
			t.Log("Test passed, panic was caught!")
		} else {
			t.Fatalf("expected panic but none occurred")
		}
	}()

	NewClient(invalid.(Dialer))
}

func TestWhenGivenAClientItShouldReturnSuccess(t *testing.T) {
	client := NewClient(&SpyDialer{})

	if client == nil {
		t.Fatalf("expected non-nil client")
	}

	if client.Dialer == nil {
		t.Fatalf("expected non-nil dialer in client")
	}
}

func TestClientConnectCallsDialerWithCorrectParameters(t *testing.T) {
	client := NewClient(spyDialer)

	urlStr := "ws://example.com/socket"
	headers := map[string][]string{
		"Authorization": {"Bearer token"},
	}

	err := client.Connect(urlStr, headers)

	// SpyDialer doesn't establish a real connection
	assert.Nil(t, err)
	assert.Equal(t, spyDialer.CalledWithURL, urlStr)
	assert.Equal(t, len(headers), len(spyDialer.CalledWithHeader))
}
