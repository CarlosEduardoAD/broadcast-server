package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func BypassCheck(r *http.Request) bool {
	return true
}

func GenerateRandomIp() string {
	random := rand.New(rand.NewSource(99))
	mockIp := "192.168.0."

	return mockIp + fmt.Sprint(random.Int63())
}

func CheckIfFileExists(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return true, nil
}
