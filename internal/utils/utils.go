package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

const (
	ALLOWED_ORIGIN = "http://localhost:8080"
)

func BypassCheck(r *http.Request) bool {
	return true
}

func Authorize(r *http.Request) bool {
	org := CheckOrigin(r)
	if !org {
		return false
	}

	auth := CheckAuthorization(r)

	if !auth {
		return false
	}

	return true
}

func CheckAuthorization(r *http.Request) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return false
	}

	encodedCredentials := parts[1]
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return false
	}

	credentials := string(decodedBytes)
	expectedCredentials := os.Getenv("BASIC_AUTH_CREDENTIALS")
	if credentials != expectedCredentials {
		return false
	}

	return true
}

func CheckOrigin(r *http.Request) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allowed := os.Getenv("ALLOWED_ORIGIN")

	origin := r.Header.Get("Origin")

	if origin != allowed {
		return false
	}

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

func Fork() (int, error) {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	// Add env to run process as daemon
	cmd.Env = append(os.Environ(), "IS_DAEMON=1")
	// Optional: redirect input/outputs
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Setsid is used to detach the process from the parent (normally a shell)
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	os.WriteFile("broadcast-server.pid", []byte(fmt.Sprint(cmd.Process.Pid)), 0644)

	return cmd.Process.Pid, nil
}
