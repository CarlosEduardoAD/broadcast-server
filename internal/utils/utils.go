package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"syscall"
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
