package healthcheck

import (
	"log"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Up and running!")
}
