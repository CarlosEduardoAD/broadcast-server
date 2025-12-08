package utils

import (
	"fmt"
	"math/rand"
	"net/http"
)

func BypassCheck(r *http.Request) bool {
	return true
}

func GenerateRandomIp() string {
	random := rand.New(rand.NewSource(99))
	mockIp := "192.168.0."

	return mockIp + fmt.Sprint(random.Int63())

}
